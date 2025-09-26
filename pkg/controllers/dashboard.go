package controllers

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func LoadDashboard(c echo.Context) error {
	var dashboardData dto.DashboardResponse

	userID := c.Get(utils.UserContextKey).(uuid.UUID)

	theUser, err := utils.Queries.GetUserById(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "Failed",
			"message": "Could not get the required user from database",
			"user_id": userID,
			"error":   err.Error(),
		})
	}

	submissions, err := utils.Queries.GetSubmissionByUser(c.Request().Context(), userID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "Failed",
			"message": "could not get submissions",
			"error":   err.Error(),
		})
	}

	current_round, err := utils.Queries.GetUserRound(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "Failed",
			"message": "Could not retrieve the user's current round",
			"error":   err.Error(),
		})
	}

	dashboardData.CurrentRound = int8(current_round)

	questionids := make(map[uuid.UUID]dto.RoundPoints)

	for _, submission := range submissions {
		question, err := utils.Queries.GetQuestion(c.Request().Context(), submission.QuestionID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":        "Failed",
				"message":       "Could not get the question for a submission",
				"submission_id": submission.ID,
				"question_id":   submission.QuestionID,
				"error":         err.Error(),
			})
		}

		if _, ok := questionids[question.ID]; !ok {
			questionids[question.ID] = dto.RoundPoints{Round: int(question.Round), Points: 0}

		}

		sub_result, err := utils.Queries.GetSubmissionResultsBySubmissionID(c.Request().Context(), submission.ID)

		if err == nil {
			for _, s := range sub_result {
				if s.PointsAwarded > int32(questionids[question.ID].Points) {
					questionids[question.ID] = dto.RoundPoints{Round: int(question.Round), Points: s.PointsAwarded}
				}
			}
		}
	}

	round_scores := [4]int32{0, 0, 0, 0}
	questions_completed := [4]int32{0, 0, 0, 0}
	questions_not_completed := [4]int32{0, 0, 0, 0}

	for _, qdata := range questionids {
		switch qdata.Round {
		case 0:
			questions_completed[0] += 1
			round_scores[0] += qdata.Points
		case 1:
			questions_completed[1] += 1
			round_scores[1] += qdata.Points
		case 2:
			questions_completed[2] += 1
			round_scores[2] += qdata.Points
		case 3:
			questions_completed[3] += 1
			round_scores[3] += qdata.Points
		}
	}

	questions_not_completed[0] = 2 - questions_completed[0]
	questions_not_completed[1] = 8 - questions_completed[1]
	questions_not_completed[2] = 7 - questions_completed[2]
	questions_not_completed[3] = 3 - questions_completed[3]

	dashboardData.UserName = theUser.Name
	dashboardData.Email = theUser.Email

	dashboardData.QuestionsCompleted = questions_completed
	dashboardData.QuestionsNotCompleted = questions_not_completed
	dashboardData.RoundScores = round_scores

	dashboardData.CurrentRound = int8(theUser.RoundQualified)

	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   dashboardData,
	})
}
