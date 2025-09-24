package dto

type DashboardResponse struct {
	UserName              string   `json:"username"`
	Email                 string   `json:"email"`
	QuestionsCompleted    [4]int32 `json:"questionsCompleted"`
	QuestionsNotCompleted [4]int32 `json:"questionnsNotCompleted"`
	RoundScores           [4]int32 `json:"roundScores"`
	CurrentRound          int8     `json:"currentRound"`
}

type RoundPoints struct {
	Round  int
	Points int32
}
