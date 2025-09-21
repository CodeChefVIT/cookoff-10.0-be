package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	submissions "github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/submission"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RunCodeResponse struct {
	Result          []submissions.Judgeresp `json:"result"`
	TestCasesPassed int                     `json:"no_testcases_passed"`
}

func RunCode(c echo.Context) error {
	var req dto.SubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid questionID"})
	}

	ctx := context.Background()
	testcasesRows, err := utils.Queries.GetTestCasesByQuestion(ctx, questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch testcases"})
	}
	if len(testcasesRows) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No testcases found for this question"})
	}

	params := url.Values{}
	params.Add("base64_encoded", "true")
	params.Add("wait", "true")

	runtimeMut, err := submissions.RuntimeMut(req.LanguageID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	response := RunCodeResponse{
		Result: make([]submissions.Judgeresp, len(testcasesRows)),
	}

	for i, tc := range testcasesRows {
		runtimeVal, err := tc.Runtime.Float64Value()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid runtime value"})
		}

		finalRuntime := runtimeVal.Float64 * float64(runtimeMut)

		payload := submissions.Submission{
			LanguageID:     req.LanguageID,
			SourceCode:     submissions.B64(req.SourceCode),
			Stdin:          submissions.B64(tc.Input),
			ExpectedOutput: submissions.B64(tc.ExpectedOutput),
			Runtime:        finalRuntime,
		}

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Error marshaling payload: %v", err),
			})
		}

		result, err := submissions.SendToJudge0(params, payloadJSON)
		if err != nil {
			fmt.Printf("Error sending request to Judge0: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error sending request to Judge0"})
		}
		defer result.Body.Close()

		bodyBytes, err := io.ReadAll(result.Body)
		if err != nil {
			fmt.Printf("Error reading Judge0 response body: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read Judge0 response"})
		}

		fmt.Println("Raw Judge0 Response:", string(bodyBytes))

		var data submissions.Judgeresp
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			fmt.Printf("Error decoding response from Judge0: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error decoding response from Judge0"})
		}

		data.Stdout, _ = submissions.DecodeB64(data.Stdout)
		data.Stderr, _ = submissions.DecodeB64(data.Stderr)
		data.Message, _ = submissions.DecodeB64(data.Message)

		response.Result[i] = data

		if data.Status.ID == 3 {
			response.TestCasesPassed++
		}
	}

	return c.JSON(http.StatusOK, response)
}
