package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	submissions "github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/submission"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Token struct {
	Token string `json:"token"`
}

func SubmitCode(c echo.Context) error {
	var req dto.SubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID, ok := c.Get(utils.UserContextKey).(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid questionID"})
	}

	allowed, err := auth.VerifyRoundAccess(c.Request().Context(), userID, questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to verify round access"})
	}
	if !allowed {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You are not qualified for this round yet"})
	}

	submissionID := uuid.New()

	ctx := c.Request().Context()

	payload, testcase_id, err := submissions.CreateSubmission(ctx, questionID, req.LanguageID, req.SourceCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create submission"})
	}

	judge0URL, err := url.Parse(utils.Config.Judge0URI + "/submissions/batch")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Submission failed",
			"error":  "Error parsing judge0 url",
		})
	}

	params := url.Values{}
	params.Add("base64_encoded", "true")

	// judge0URL.RawQuery = params.Encode()
	// resp, err := http.Post(judge0URL.String(), "application/json", bytes.NewBuffer(payload))
	resp, err := submissions.SendToJudge(judge0URL, params, payload)

	if err != nil {
		logger.Errorf("Error sending request to Judge0: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error sending request",
		})
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error reading response body from Judge0: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error reading response body",
		})
	}

	if resp.StatusCode != http.StatusCreated {
		logger.Errorf("Judge0 returned status code %d", resp.StatusCode)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error from Judge0",
		})
	}

	var tokens []Token
	_ = json.Unmarshal(respBytes, &tokens)

	for i, t := range tokens {
		err := utils.TokenCache.Set(ctx, fmt.Sprintf("token:%s", t.Token), fmt.Sprintf("%s:%s", submissionID, testcase_id[i]), 0).Err()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to cache token",
			})
		}

		err = utils.TokenCache.SAdd(ctx, fmt.Sprintf("sub:%s:tokens", submissionID), t.Token).Err()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to add token to set",
			})
		}
	}

	err = utils.Queries.CreateSubmission(ctx, db.CreateSubmissionParams{
		ID:         submissionID,
		UserID:     userID,
		QuestionID: questionID,
		LanguageID: int32(req.LanguageID),
		SourceCode: req.SourceCode,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to create batch submission in database",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"submission_id": submissionID,
	})
}

// GetUserSubmissions fetches all submissions of a user with their results
func GetUserSubmissions(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := utils.Queries.GetUserById(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch user",
		})
	}

	subs, err := utils.Queries.GetSubmissionsByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch submissions",
		})
	}

	var submissionsWithResults []dto.SubmissionWithResults

	for _, sub := range subs {
		results, err := utils.Queries.GetSubmissionResultsBySubmissionID(c.Request().Context(), sub.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": fmt.Sprintf("Failed to fetch results for submission %s", sub.ID),
			})
		}

		submission := db.Submission{
			ID:              sub.ID,
			UserID:          sub.UserID,
			QuestionID:      sub.QuestionID,
			LanguageID:      sub.LanguageID,
			SourceCode:      sub.SourceCode,
			TestcasesPassed: sub.TestcasesPassed,
			TestcasesFailed: sub.TestcasesFailed,
			Runtime:         sub.Runtime,
			Memory:          sub.Memory,
			Status:          sub.Status,
			Description:     sub.Description,
		}

		submissionsWithResults = append(submissionsWithResults, dto.SubmissionWithResults{
			Submission: submission,
			Results:    results,
		})
	}

	response := dto.UserSubmissionsResponse{
		User:        user,
		Submissions: submissionsWithResults,
	}

	return c.JSON(http.StatusOK, response)
}
