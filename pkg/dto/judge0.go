package dto

// Judge0CallbackPayload represents the payload for Judge0 callback processing
type Judge0CallbackPayload struct {
	SubmissionID string `json:"submission_id"`
	Token        string `json:"token"`
	Status       struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
	Language struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"language"`
	Stdin         string `json:"stdin"`
	Source        string `json:"source"`
	Output        string `json:"output"`
	Stderr        string `json:"stderr"`
	CompileOutput string `json:"compile_output"`
	Runtime       int    `json:"runtime"`
	Memory        int    `json:"memory"`
	CreatedAt     string `json:"created_at"`
	FinishedAt    string `json:"finished_at"`
}
