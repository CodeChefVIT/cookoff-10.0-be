package dto

type DashboardResponse struct {
	UserName              string   `json:"username"`
	Email                 string   `json:"email"`
	QuestionsCompleted    [4]int32 `json:"questions_completed"`
	QuestionsNotCompleted [4]int32 `json:"questions_not_completed"`
	RoundScores           [4]int32 `json:"round_scores"`
	CurrentRound          int8     `json:"current_round"`
}

type RoundPoints struct {
	Round  int
	Points int32
}
