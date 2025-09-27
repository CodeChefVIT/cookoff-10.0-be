package dto

type CreateTestCaseRequest struct {
	ExpectedOutput string `json:"expected_output" validate:"required"`
	Memory         string `json:"memory" validate:"required"`
	Input          string `json:"input" validate:"required"`
	Hidden         bool   `json:"hidden"`
	Runtime        string `json:"runtime" validate:"required"`
	QuestionID     string `json:"question_id" validate:"required,uuid"`
}

type UpdateTestCaseRequest struct {
	ExpectedOutput string      `json:"expected_output"`
	Memory         string 		`json:"memory"`
	Input          string      `json:"input"`
	Hidden         *bool       `json:"hidden"`
	Runtime        string 		`json:"runtime"`
	QuestionID     string      `json:"question_id" validate:"omitempty,uuid"`
}
