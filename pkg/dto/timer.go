package dto

type SetTimeRequest struct {
	RoundID string `json:"round_id" validate:"required"`
	Time    string `json:"time" validate:"reqired"`
}
