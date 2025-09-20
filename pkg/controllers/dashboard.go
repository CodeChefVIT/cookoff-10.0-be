package controllers

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func LoadDashboard(c echo.Context) error {
	var dashboardData dto.DashboardResponse

	userID, err := auth.GetUserID(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "Failed",
			"message": "Could not get user_id",
			"error":   err.Error(),
		})
	}

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

	var r0s, r1s, r2s, r3s = 0, 0, 0, 0
	var qr0c, qr1c, qr2c, qr3c = 0, 0, 0, 0

	for _, qdata := range questionids {
		switch qdata.Round {
		case 0:
			qr0c += 1
			r0s += int(qdata.Points)
		case 1:
			qr1c += 1
			r1s += int(qdata.Points)
		case 2:
			qr2c += 1
			r2s += int(qdata.Points)
		case 3:
			qr3c += 1
			r3s += int(qdata.Points)
		}
	}

	dashboardData.UserName = theUser.Name
	dashboardData.Email = theUser.Email

	dashboardData.QuestionsCompleted0 = int64(qr0c)
	dashboardData.QuestionsCompleted1 = int64(qr1c)
	dashboardData.QuestionsCompleted2 = int64(qr2c)
	dashboardData.QuestionsCompleted3 = int64(qr3c)

	dashboardData.QuestionsNotCompleted0 = int64(2 - qr0c)
	dashboardData.QuestionsNotCompleted1 = int64(8 - qr1c)
	dashboardData.QuestionsNotCompleted2 = int64(7 - qr2c)
	dashboardData.QuestionsNotCompleted3 = int64(3 - qr3c)

	dashboardData.Round0Score = int32(r0s)
	dashboardData.Round1Score = int32(r1s)
	dashboardData.Round2Score = int32(r2s)
	dashboardData.Round3Score = int32(r3s)

	dashboardData.CurrentRound = int8(theUser.RoundQualified)

	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"data":   dashboardData,
	})
}
