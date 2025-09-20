package controllers

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateQuestion(c echo.Context) error {
	var req db.CreateQuestionParams
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not create question",
			"error":  err.Error(),
		})
	}
	req.ID = uuid.New()
	if err := utils.Queries.CreateQuestion(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not create question",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"status":  "success",
		"message": "question created",
		"data":    req,
	})
}

func GetQuestion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Question not found!",
			"error":  "UUID GALAT HAI BHAI",
		})
	}
	q, err := utils.Queries.GetQuestion(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status": "Could not get question",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":   "success",
		"question": q,
	})
}

func GetAllQuestions(c echo.Context) error {
	questions, err := utils.Queries.GetAllQuestions(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not get all the questions",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":    "success",
		"questions": questions,
	})
}

func GetQuestionsByRound(c echo.Context) error {

	userID, err := auth.GetUserID(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "Failed",
			"message": "Error Getting user_id",
			"error":   err.Error(),
		})
	}

	round, err := utils.Queries.GetUserRound(c.Request().Context(), userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "Failed",
			"message": "could not get the users current round",
			"error":   err.Error(),
		})
	}

	if round < 0 || round > 3 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  "Failed",
			"message": "round number is invalid",
		})
	}

	questions, err := utils.Queries.GetQuestionsByRound(c.Request().Context(), int32(round))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "Failed",
			"message": "unable to fetch the questions",
			"error":   err.Error(),
		})
	}

	result := []echo.Map{}

	for _, q := range questions {
		testcases, err := utils.Queries.GetTestCasesByQuestion(c.Request().Context(), q.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":      "Failed",
				"message":     "Could not get the testcases for question",
				"question_id": q.ID,
				"error":       err.Error(),
			})
		}
		result = append(result, echo.Map{
			"question":  q,
			"testcases": testcases,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":              "success",
		"round":               round,
		"questions_testcases": result,
	})
}

func UpdateQuestion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not update question",
			"error":  "UUID GALAT HAI BHAI",
		})
	}
	var req db.UpdateQuestionParams
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not update question",
			"error":  err.Error(),
		})
	}
	req.ID = id
	if err := utils.Queries.UpdateQuestion(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not update question",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "Question updated successfully!",
	})
}

func DeleteQuestion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not delete question",
			"error":  "UUID GALAT HAI BHAI",
		})
	}
	if err := utils.Queries.DeleteQuestion(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not delete question",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "Question deleted successfully",
	})
}

func ActivateBounty(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not activate bounty for the question",
			"error":  "UUID GALAT HAI BHAI",
		})
	}
	if err := utils.Queries.UpdateQuestionBountyActive(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not activate bounty for the question",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "Bounty activated for the question",
	})
}

func DeactivateBounty(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not deactivate bounty for the question",
			"error":  "UUID GALAT HAI BHAI",
		})
	}
	if err := utils.Queries.UpdateQuestionBountyInactive(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not deactivate bounty for the question",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "Bounty deactivated for the question",
	})
}
