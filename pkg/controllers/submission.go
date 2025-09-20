package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	submissions "github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/submission"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"

	//"github.com/CodeChefVIT/cookoff-10.0-be/pkg/workers"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var Queries *db.Queries

func SubmitCode(c echo.Context) error {
	var req dto.SubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	submissionID := uuid.New()

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid questionID"})
	}

	ctx := context.Background()
	testcasesRows, err := Queries.GetTestCasesByQuestion(ctx, questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch testcases"})
	}

	var testcases []map[string]string
	for _, tc := range testcasesRows {
		testcases = append(testcases, map[string]string{
			"input":  tc.Input,
			"output": tc.ExpectedOutput,
		})
	}

	tokens, err := submissions.CreateBatchSubmission(submissionID.String(), req.SourceCode, req.LanguageID, testcases)
	if err != nil {
		fmt.Println("CreateBatchSubmission error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create batch submission"})
	}
	for _, token := range tokens {
        if err := utils.TokenCache.Set(utils.Ctx, token, submissionID, 0).Err(); err != nil {
            fmt.Printf("Failed to cache token %s: %v\n", token, err)
        }
    }
	sub := utils.SubmissionInput{
		ID:         submissionID,
		QuestionID: req.QuestionID,
		LanguageID: req.LanguageID,
		SourceCode: req.SourceCode,
		UserID:     req.UserID,
	}
	if err := utils.SaveSubmission(sub); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save submission record"})
	}


	return c.JSON(http.StatusOK, map[string]interface{}{
		"submission_id": submissionID,
	})
}