package dto

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
)

type SubmissionRequest struct {
	SourceCode string `json:"source_code" validate:"required"`
	LanguageID int    `json:"language_id" validate:"required"`
	QuestionID string `json:"question_id" validate:"required"`
}
type SubmissionPayload struct {
	LanguageID int     `json:"language_id"`
	SourceCode string  `json:"source_code"`
	Input      string  `json:"stdin"`
	Output     string  `json:"expected_output"`
	Runtime    float64 `json:"cpu_time_limit"`
}
type CustomSubmissionRequest struct {
	SourceCode string `json:"source_code" validate:"required"`
	LanguageID int    `json:"language_id" validate:"required"`
	Input      string `json:"input"`
}

type UserSubmissionsResponse struct {
	User        db.User       `json:"user"`
	Submissions []SubmissionWithResults `json:"submissions"`
}

type SubmissionWithResults struct {
	Submission db.Submission        `json:"submission"`
	Results    []db.SubmissionResult `json:"results"`
}