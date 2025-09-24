package controllers

import (
	"net/http"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/labstack/echo/v4"
)

func GetTime(c echo.Context) error {
	roundID := c.Param("id")
	timeStr, err := utils.RedisClient.HGet(c.Request().Context(), roundID, "round_end_time").Result()
	if err != nil {
		logger.Errorf("Round get time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	roundEndTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"round_end_time": roundEndTime,
		"server_time":    time.Now(),
	})
}

func SetTime(c echo.Context) error {
	var req dto.SetTimeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "failed",
			"error":  "invalid request",
		})
	}

	if err := validator.ValidatePayload(req); err != nil {
		return c.JSON(http.StatusNotAcceptable, echo.Map{
			"status": "failed",
			"error":  "invalid input",
		})
	}

	_, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "incorrect time format",
		})
	}

	err = utils.RedisClient.HSet(c.Request().Context(), req.RoundID, "round_end_time", req.Time).Err()
	if err != nil {
		logger.Errorf("Round get time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "time set successfully",
	})
}

func UpdateTime(c echo.Context) error {
	return nil
}
