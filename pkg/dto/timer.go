package dto

type SetTimeRequest struct {
	RoundID string `json:"round_id" validate:"required"`
	Time    string `json:"time" validate:"required"`
}

type UpdateTimeRequest struct {
	RoundID  string `json:"round_id" validate:"required"`
	Duration string `json:"duration" validate:"required"`
}
