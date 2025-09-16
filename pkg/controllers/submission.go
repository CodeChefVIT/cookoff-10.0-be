package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	submissions "github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/submission"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

var Queries *db.Queries
var Rdb *redis.Client

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


	for _, t := range tokens {
		Rdb.Set(ctx, "token:"+t, submissionID.String(), 0)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "submission successful",
	})
}
