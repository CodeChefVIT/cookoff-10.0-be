package submissions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type SubmissionInput struct {
	ID         uuid.UUID
	QuestionID string
	LanguageID int
	SourceCode string
	UserID     string
}

type Judge0Submission struct {
	LanguageID     int     `json:"language_id"`
	SourceCode     string  `json:"source_code"`
	Stdin          string  `json:"stdin,omitempty"`
	ExpectedOutput string  `json:"expected_output,omitempty"`
	Runtime        float64 `json:"cpu_time_limit"`
	Callback       string  `json:"callback_url,omitempty"`
}

func CreateBatchSubmission(submissionID, sourceCode string, languageID int, testcases []map[string]string) ([]string, error) {
	var submissions []Judge0Submission

	callbackURL := os.Getenv("CALLBACK_URL")
	if callbackURL == "" {
		return nil, errors.New("CALLBACK_URL not set in environment")
	}

	runtime_mut, err := RuntimeMut(languageID)
	if err != nil {
		return nil, err
	}

	for _, tc := range testcases {
		submissions = append(submissions, Judge0Submission{
			LanguageID:     languageID,
			SourceCode:     sourceCode,
			Stdin:          tc["input"],
			ExpectedOutput: tc["output"],
			Runtime:        float64(runtime_mut),
			Callback:       callbackURL,
		})
	}

	if len(submissions) == 0 {
		return nil, errors.New("no testcases provided for batch submission")
	}

	payload := map[string]interface{}{"submissions": submissions}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal submissions: %v", err)
	}

	judge0URI := os.Getenv("JUDGE0_URI")
	if judge0URI == "" {
		return nil, errors.New("JUDGE0_URI not set in environment")
	}


	batchURL := judge0URI + "/submissions/batch?base64_encoded=true"

	req, err := http.NewRequest("POST", batchURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-RapidAPI-Key", apiKey)
	// req.Header.Set("X-RapidAPI-Host", "judge0-ce.p.rapidapi.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var respData []struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(bodyBytes, &respData); err != nil {
		return nil, fmt.Errorf("failed to decode response JSON: %v", err)
	}

	if len(respData) == 0 {
		return nil, errors.New("no tokens returned from Judge0")
	}

	tokens := make([]string, len(respData))
	for i, t := range respData {
		tokens[i] = t.Token
	}

	return tokens, nil
}
