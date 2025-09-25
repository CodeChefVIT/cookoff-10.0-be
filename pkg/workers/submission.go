package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	TypeProcessSubmission = "submission:process"
	SubmissionDoneStatus  = "DONE"
)

// ProcessJudge0CallbackTask simply logs the received Judge0 callback
func ProcessJudge0CallbackTask(ctx context.Context, t *asynq.Task) error {
	var data dto.Judge0CallbackPayload

	logger.Infof("Processing task: %v", t.Type)
	logger.Infof("Payload: %v", string(t.Payload()))

	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		log.Printf("Error unmarshalling task payload: %v\n", err)
		return err
	}

	timeValue, err := parseTime(data.Time)
	if err != nil {
		log.Println("Error parsing time value: ", err)
		return err
	}

	value, testcase, err := utils.GetSubmissionIDByToken(ctx, data.Token)
	if err != nil {
		log.Println("Error getting submission ID from token: ", err)
		return err
	}

	idUUID, err := uuid.Parse(value)
	if err != nil {
		log.Fatalf("Error parsing UUID: %v", err)
	}

	testidUUID, err := uuid.Parse(testcase)
	if err != nil {
		log.Fatalf("Error parsing UUID: %v", err)
	}

	var status string
	switch data.Status.ID {
	case "1":
		status = "In Queue"
	case "2":
		status = "Processing"
	case "3":
		status = "success"
	case "4":
		status = "wrong answer"
	case "5":
		status = "Time Limit Exceeded"
	case "6":
		status = "Compilation error"
	case "7", "8", "9", "10", "11", "12":
		status = "Runtime error"
	case "13":
		status = "Internal Error"
	case "14":
		status = "Exec Format Error"
	}

	if status != "" {
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			status,
		)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	if err := utils.DeleteToken(ctx, data.Token); err != nil {
		log.Println("Error deleting token: ", err)
		return err
	}

	tokenCount, err := utils.GetTokenCount(ctx, value)
	if err != nil {
		log.Println("Error getting token count: ", err)
		return err
	}

	fmt.Println("Token :- ", tokenCount)

	if tokenCount == 0 {
		err = utils.UpdateSubmission(ctx, idUUID)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseTime(timeStr string) (float64, error) {
	if timeStr == "" {
		log.Println("Time value is empty, setting time to 0 for this submission.")
		return 0, nil
	}

	timeValue, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		return 0, err
	}
	return timeValue, nil
}

func handleCompilationError(
	ctx context.Context,
	idUUID uuid.UUID,
	data dto.Judge0CallbackPayload,
	time int,
	testcase uuid.UUID,
	status string,
) error {
	subID, err := uuid.NewV7()
	if err != nil {
		log.Println("Error updating submission for compilation error: ", err)
		return err
	}

	err = utils.Queries.CreateSubmissionResult(ctx, db.CreateSubmissionResultParams{
		ID:            subID,
		TestcaseID:    uuid.NullUUID{UUID: testcase, Valid: true},
		SubmissionID:  idUUID,
		Runtime:       pgtype.Numeric{Int: big.NewInt(int64(time)), Valid: true},
		Memory:        pgtype.Numeric{Int: big.NewInt(int64(data.Memory)), Valid: true},
		PointsAwarded: 10,
		Status:        status,
		Description:   &data.Status.Description,
	})
	if err != nil {
		log.Println("Error creating submission status error: ", err)
		return err
	}
	return nil
}
