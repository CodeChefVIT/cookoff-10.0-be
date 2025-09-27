package dto

import "github.com/jackc/pgx/v5/pgtype"

type CreateTestCaseRequest struct {
	ExpectedOutput string         `json:"expected_output" validate:"required"`
	Memory         pgtype.Numeric `json:"memory" validate:"required"`
	Input          string         `json:"input" validate:"required"`
	Hidden         bool           `json:"hidden"`
	Runtime        pgtype.Numeric `json:"runtime" validate:"required"`
	QuestionID     string         `json:"question_id" validate:"required,uuid"`
}

type UpdateTestCaseRequest struct {
	ExpectedOutput string         `json:"expected_output"`
	Memory         pgtype.Numeric `json:"memory"`
	Input          string         `json:"input"`
	Hidden         *bool          `json:"hidden"`
	Runtime        pgtype.Numeric `json:"runtime"`
	QuestionID     string         `json:"question_id" validate:"omitempty,uuid"`
}
