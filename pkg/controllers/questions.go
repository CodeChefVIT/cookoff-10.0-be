package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateQuestion(c echo.Context) error {

	var req dto.CreateQuestion
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not create question",
			"error":  err.Error(),
		})
	}
	if err := utils.Queries.CreateQuestion(c.Request().Context(), db.CreateQuestionParams{
		ID:               uuid.New(),
		Description:      req.Description,
		Title:            req.Title,
		Qtype:            req.Qtype,
		Isbountyactive:   req.Isbountyactive,
		InputFormat:      req.InputFormat,
		Points:           req.Points,
		Round:            req.Round,
		Constraints:      req.Constraints,
		OutputFormat:     req.OutputFormat,
		SampleTestInput:  req.SampleTestInput,
		SampleTestOutput: req.SampleTestOutput,
		Explanation:      req.Explanation,
	}); err != nil {
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
	userID := c.Get(utils.UserContextKey).(uuid.UUID)

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

	cache_round_key := fmt.Sprintf("questions_round_%d", round)
	qs, err := utils.RedisClient.Get(c.Request().Context(), cache_round_key).Result()

	if err == nil {
		result := []map[string]any{}
		err := json.Unmarshal([]byte(qs), &result)
		if err == nil {
			return c.JSON(http.StatusOK, echo.Map{
				"status":              "success",
				"round":               round,
				"questions_testcases": result,
			})
		}
	}

	questions, err := utils.Queries.GetQuestionsByRound(c.Request().Context(), int32(round))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status":  "Failed",
			"message": "unable to fetch the questions",
			"error":   err.Error(),
		})
	}

	result := []map[string]any{}

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
		result = append(result, map[string]any{
			"question":  q,
			"testcases": testcases,
		})
	}

	qdata, _ := json.Marshal(result)
	utils.RedisClient.Set(c.Request().Context(), cache_round_key, qdata, 2*time.Minute)

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
	var req dto.Question
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not update question",
			"error":  err.Error(),
		})
	}
	if err := utils.Queries.UpdateQuestion(c.Request().Context(), db.UpdateQuestionParams{
		Description:      req.Description,
		Title:            req.Title,
		Qtype:            req.Qtype,
		Isbountyactive:   req.Isbountyactive,
		InputFormat:      req.InputFormat,
		Points:           req.Points,
		Round:            req.Round,
		Constraints:      req.Constraints,
		OutputFormat:     req.OutputFormat,
		SampleTestInput:  req.SampleTestInput,
		SampleTestOutput: req.SampleTestOutput,
		Explanation:      req.Explanation,
		ID:               id,
	}); err != nil {
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
