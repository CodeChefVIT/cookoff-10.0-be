package dto

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
