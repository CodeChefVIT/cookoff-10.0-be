package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	submissions "github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/submission"
	"github.com/labstack/echo/v4"
)

type RunCustomCodeResponse struct {
	Result submissions.Judgeresp `json:"result"`
}

func RunCustom(c echo.Context) error {
	var req dto.CustomSubmissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	params := url.Values{}
	params.Add("base64_encoded", "true")
	params.Add("wait", "true")

	
	payload := submissions.Submission{
		LanguageID: req.LanguageID,
		SourceCode: submissions.B64(req.SourceCode),
		Stdin:      submissions.B64(req.Input),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Error marshaling payload: %v", err)})
	}

	result, err := submissions.SendToJudge0(params, payloadJSON)
	if err != nil {
		fmt.Printf("Error sending request to Judge0: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error sending request to Judge0"})
	}
	defer result.Body.Close()

	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read Judge0 response"})
	}

	fmt.Println("Raw Judge0 Response:", string(bodyBytes))

	var data submissions.Judgeresp
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error decoding response from Judge0"})
	}

	
	data.Stdout, _ = submissions.DecodeB64(data.Stdout)
	data.Stderr, _ = submissions.DecodeB64(data.Stderr)
	data.Message, _ = submissions.DecodeB64(data.Message)

	response := RunCustomCodeResponse{Result: data}
	return c.JSON(http.StatusOK, response)
}
