package workers

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"strings"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

// ProcessJudge0CallbackTask simply logs the received Judge0 callback
func ProcessJudge0CallbackTask(ctx context.Context, t *asynq.Task) error {
	var data dto.Judge0CallbackPayload
	log.Printf("Received Judge0 callback task: %s", string(t.Payload()))

	logger.Infof("Processing task: %v", t.Type)
	logger.Infof("Payload: %v", string(t.Payload()))

	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		log.Printf("Error unmarshalling task payload: %v\n", err)
		return err
	}

	timeValue, err := parseTime(data.Runtime)
	if err != nil {
		log.Println("Error parsing time value: ", err)
		return err
	}

	subID, err := utils.GetSubmissionIDByToken(data.Token)
	if err != nil {
		log.Println("Error getting submission ID from token: ", err)
		return err
	}

	temp := strings.Split(subID, ":")
	value := temp[0]
	testcase := temp[1]

	idUUID, err := uuid.Parse(value)
	if err != nil {
		log.Fatalf("Error parsing UUID: %v", err)
	}

	testidUUID, err := uuid.Parse(testcase)
	if err != nil {
		log.Fatalf("Error parsing UUID: %v", err)
	}

	switch data.Status.ID {
	case 1:
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"In Queue",
		)
	case 2:
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Processing",
		)
	case 3:
		// testcasesPassed++
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"success")
	case 4:
		// testcasesFailed++
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"wrong answer",
		)
	case 5:
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Time Limit Exceeded",
		)
	case 6:
		// testcasesFailed++
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Compilation error",
		)
	case 7, 8, 9, 10, 11, 12:
		// testcasesFailed++
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Runtime error",
		)
	case 13:
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Internal Error",
		)
	case 14:
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Exec Format Error",
		)
	}

	if err != nil {
		return err
	}

	if err := utils.DeleteTokensBySubmissionID(subID); err != nil {
		log.Println("Error deleting token: ", err)
		return err
	}

	// tokenCount, err := submission.Tokens.GetTokenCount(ctx, value)
	// if err != nil {
	// 	log.Println("Error getting token count: ", err)
	// 	return err
	// }

	// fmt.Println("Token :- ", tokenCount)

	// if tokenCount == 0 {
	// 	err = submission.UpdateSubmission(ctx, idUUID)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func parseTime(time int) (float64, error) {
	timeValue := float64(time)
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
