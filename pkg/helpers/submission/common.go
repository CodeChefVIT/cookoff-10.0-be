package submissions

import (
	"bytes"
	"encoding/base64"
	//"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Judgeresp struct {
	Token    string `json:"token"`
	Status   Status `json:"status"`
	Stdout   string `json:"stdout"`
	Time     string `json:"time"`
	Memory   int    `json:"memory"`
	Stderr   string `json:"stderr"`
	Message  string `json:"message"`
	Language string `json:"language"`
}

type Status struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Submission struct {
	LanguageID     int     `json:"language_id"`
	SourceCode     string  `json:"source_code"`
	Stdin          string  `json:"stdin,omitempty"`
	ExpectedOutput string  `json:"expected_output,omitempty"`
	Runtime        float64 `json:"cpu_time_limit"`
	Callback       string  `json:"callback_url,omitempty"`
}

func RuntimeMut(languageID int) (int, error) {
	switch languageID {
	case 50, 54, 60, 73, 63:
		return 1, nil
	case 51, 62:
		return 2, nil
	case 68:
		return 3, nil
	case 71:
		return 5, nil
	default:
		return 0, fmt.Errorf("invalid language ID: %d", languageID)
	}
}

func B64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func DecodeB64(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	return string(decoded), err
}
func SendToJudge0(params url.Values, payload []byte) (*http.Response, error) {
	baseURI := os.Getenv("JUDGE0_URI")
	if baseURI == "" {
		return nil, errors.New("JUDGE0_URI not set in environment")
	}

	// Make sure baseURI doesn't have a trailing slash
	if baseURI[len(baseURI)-1] == '/' {
		baseURI = baseURI[:len(baseURI)-1]
	}

	// Append /submissions explicitly
	fullURL := baseURI + "/submissions?" + params.Encode()

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Optional: Add API key if using RapidAPI
	// req.Header.Set("X-RapidAPI-Key", os.Getenv("RAPIDAPI_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Judge0: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("⚠️ Judge0 returned status %d: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("Judge0 returned non-200 status: %d", resp.StatusCode)
	}

	return resp, nil
}