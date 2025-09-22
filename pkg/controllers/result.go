package controllers

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/submission"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type ResultReq struct {
	SubID string `json:"submission_id" validate:"required"`
}

func GetResult(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Minute)
	defer cancel()

	subIDStr := c.Param("submission_id")
	subID, err := uuid.Parse(subIDStr)
	if err != nil {
		return c.JSON(400, echo.Map{"error": "Invalid UUID format"})
	}
	processed, err := submissions.CheckStatus(ctx, subID)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "Internal server error while getting status"})
	}

	if processed {
		result, err := submissions.GetSubResult(ctx, subID)
		if err != nil {
			return c.JSON(500, echo.Map{"error": "Internal server error while getting submission result"})
		}
		return c.JSON(200, result)
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return c.JSON(408, echo.Map{"error": "Submission not processed"})

		case <-ticker.C:
			processed, err := submissions.CheckStatus(ctx, subID)
			if err != nil {
				return c.JSON(500, echo.Map{"error": "Internal server error while getting status"})
			}

			if processed {
				result, err := submissions.GetSubResult(ctx, subID)
				if err != nil {
					return c.JSON(500, echo.Map{"error": "Internal server error while getting submission result"})
				}
				return c.JSON(200, result)
			}
		}
	}
}

type GetStatus struct {
	Description string `json:"description"`
	ID          int    `json:"id"`
}

type GetSub struct {
	CompileOutput *string   `json:"compile_output"`
	Memory        *int      `json:"memory"`
	Message       *string   `json:"message"`
	Status        GetStatus `json:"status"`
	Stderr        *string   `json:"stderr"`
	Stdout        *string   `json:"stdout"`
	Time          *string   `json:"time"`
	Token         string    `json:"token"`
}

type Response struct {
	Submissions []GetSub `json:"submissions"`
}

func parseTime(timeStr string) (float64, error) {
	if timeStr == "" {
		log.Println("Time value is empty, setting time to 0 for this submission.")
		return 0, nil
	}
	return strconv.ParseFloat(timeStr, 64)
}

func HandleCompilationError(
	ctx context.Context,
	idUUID uuid.UUID,
	data GetSub,
	time int,
	testcase uuid.UUID,
	status string,
) error {
	subID, err := uuid.NewV7()
	if err != nil {
		log.Println("Error generating UUID for submission status:", err)
		return err
	}

	err = utils.Queries.CreateSubmissionResult(ctx,  db.CreateSubmissionResultParams{
		ID:           subID,
		SubmissionID: idUUID,
		TestcaseID:   uuid.NullUUID{UUID: testcase, Valid: true},
		Runtime:      pgtype.Numeric{Int: big.NewInt(int64(time)), Valid: true},
		Memory:       pgtype.Numeric{Int: big.NewInt(int64(*data.Memory)), Valid: true},
		Description:  &data.Status.Description,
		Status:       status,
	})
	if err != nil {
		log.Println("Error creating submission status:", err)
		return err
	}
	return nil
}
