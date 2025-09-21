package dto

import "encoding/json"

// Judge0CallbackPayload represents the payload for Judge0 callback processing
type Judge0CallbackPayload struct {
	StdOut  *string `json:"stdout"`
	Time    string  `json:"time"`
	Memory  int     `json:"memory"`
	StdErr  *string `json:"stderr"`
	Token   string  `json:"token"`
	Message *string `json:"message"`
	Status  Status  `json:"status"`
}

type Status struct {
	ID          json.Number `json:"id"`
	Description string      `json:"description"`
}
type Judge0Response struct {
	Status struct {
		ID int `json:"id"`
	} `json:"status"`
	StdOut         string `json:"stdout"`
	StdErr         string `json:"stderr"`
	CompilerOutput string `json:"compile_output"`
	Message        string `json:"message"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	TestCaseID     string `json:"test_case_id"`
}
