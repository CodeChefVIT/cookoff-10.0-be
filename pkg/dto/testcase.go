package dto

type CreateTestCaseRequest struct {
	ExpectedOutput string      `json:"expected_output" validate:"required"`
	Memory         interface{} `json:"memory" validate:"required"`
	Input          string      `json:"input" validate:"required"`
	Hidden         bool        `json:"hidden"`
	Runtime        interface{} `json:"runtime" validate:"required"`
	QuestionID     string      `json:"question_id" validate:"required,uuid"`
}

type UpdateTestCaseRequest struct {
	ExpectedOutput string      `json:"expected_output"`
	Memory         interface{} `json:"memory"`
	Input          string      `json:"input"`
	Hidden         *bool       `json:"hidden"`
	Runtime        interface{} `json:"runtime"`
	QuestionID     string      `json:"question_id" validate:"omitempty,uuid"`
}
