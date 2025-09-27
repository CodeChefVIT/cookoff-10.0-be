package submissions

import (
	"bytes"
	"encoding/base64"

	//"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var bearer = os.Getenv("JUDGE0_TOKEN")

type Judgeresp struct {
	Token         string `json:"token"`
	Status        Status `json:"status"`
	Stdout        string `json:"stdout"`
	Time          string `json:"time"`
	Memory        int    `json:"memory"`
	Stderr        string `json:"stderr"`
	Message       string `json:"message"`
	Language      string `json:"language"`
	CompileOutput string `json:"compile_output"`
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

	if baseURI[len(baseURI)-1] == '/' {
		baseURI = baseURI[:len(baseURI)-1]
	}

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

	return resp, nil
}

func SendToJudge(judge0Url *url.URL, params url.Values, payload []byte) (*http.Response, error) {
	judge0Url.RawQuery = params.Encode()
	judgereq, err := http.NewRequest("POST", judge0Url.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request to Judge0: %v", err)
	}

	judgereq.Header.Add("Content-Type", "application/json")
	judgereq.Header.Add("Accept", "application/json")
	judgereq.Header.Add("Authorization", fmt.Sprintf("Bearer %v", bearer))

	resp, err := http.DefaultClient.Do(judgereq)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Judge0: %v", err)
	}

	return resp, nil
}
