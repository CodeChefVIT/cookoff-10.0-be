package controllers

import (
	"context"
	"log"
	"strconv"
	"time"

	submissions "github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/submission"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResultReq struct {
	SubID string `json:"submission_id" validate:"required"`
}
type ResultResp struct {
	ID             string     `json:"id"`
	QuestionID     string     `json:"question_id"`
	Passed         int        `json:"passed"`
	Failed         int        `json:"failed"`
	Runtime        float64    `json:"runtime"`
	Memory         float64    `json:"memory"`
	SubmissionTime string     `json:"submission_time"`
	Description    string     `json:"description"`
	Testcases      []TCResult `json:"testcases"`
}

type TCResult struct {
	ID             string  `json:"id"`
	Runtime        float64 `json:"runtime"`
	Memory         float64 `json:"memory"`
	Status         string  `json:"status"`
	Description    string  `json:"description"`
	ExpectedOutput string  `json:"expected_output"`
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
		return fetchResultWithTestcases(ctx, subID, c)
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
				return fetchResultWithTestcases(ctx, subID, c)
			}
		}
	}
}

func fetchResultWithTestcases(ctx context.Context, subID uuid.UUID, c echo.Context) error {
	result, err := submissions.GetSubResult(ctx, subID)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "Internal server error while getting submission result"})
	}

	questionID, err := uuid.Parse(result.QuestionID)
	if err != nil {
		return c.JSON(400, echo.Map{"error": "Invalid QuestionID UUID"})
	}

	testcases, err := utils.Queries.GetTestCasesByQuestion(ctx, questionID)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "Failed to fetch testcases"})
	}

	for i := range result.Testcases {
		if i < len(testcases) {
			result.Testcases[i].ExpectedOutput = testcases[i].ExpectedOutput
		}
	}

	return c.JSON(200, result)
}

func parseTime(timeStr string) (float64, error) {
	if timeStr == "" {
		log.Println("Time value is empty, setting time to 0 for this submission.")
		return 0, nil
	}
	return strconv.ParseFloat(timeStr, 64)
}
