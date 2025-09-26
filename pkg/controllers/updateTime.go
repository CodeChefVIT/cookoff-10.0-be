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

func UpdateTime(c echo.Context) error {
	var req dto.UpdateTimeRequest
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

	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "incorrect duration format",
		})
	}

	roundEndTimeStr, err := utils.RedisClient.HGet(c.Request().Context(), "round_end_time", req.RoundID).Result()
	if err != nil {
		logger.Errorf("Round get time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	roundEndTime, err := time.ParseInLocation(time.RFC3339, roundEndTimeStr, utils.IST)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "incorrect time format",
		})
	}

	updatedTime := roundEndTime.Add(duration)

	err = utils.RedisClient.HSet(c.Request().Context(), "round_end_time", req.RoundID, updatedTime).Err()
	if err != nil {
		logger.Errorf("Round get time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	err = utils.RedisClient.ExpireAt(c.Request().Context(), "is_round_started", updatedTime).Err()
	if err != nil {
		logger.Errorf("current_round TTL set error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "time updated successfully",
	})
}
