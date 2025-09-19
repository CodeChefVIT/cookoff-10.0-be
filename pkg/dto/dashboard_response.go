package dto

type DashboardResponse struct {
	UserName               string `json:"user_name"`
	Email                  string `json:"email"`
	QuestionsCompleted0    int64  `json:"questions_completed_round0"`
	QuestionsNotCompleted0 int64  `json:"questions_not_completed_round0"`
	QuestionsCompleted1    int64  `json:"questions_completed_round1"`
	QuestionsNotCompleted1 int64  `json:"questions_not_completed_round1"`
	QuestionsCompleted2    int64  `json:"questions_completed_round2"`
	QuestionsNotCompleted2 int64  `json:"questions_not_completed_round2"`
	QuestionsCompleted3    int64  `json:"questions_completed_round3"`
	QuestionsNotCompleted3 int64  `json:"questions_not_completed_round3"`
	Round0Score            int32  `json:"round_0_score"`
	Round1Score            int32  `json:"round_1_score"`
	Round2Score            int32  `json:"round_2_score"`
	Round3Score            int32  `json:"round_3_score"`
	CurrentRound           int8   `json:"current_round"`
}

type RoundPoints struct {
	Round  int
	Points int32
}
