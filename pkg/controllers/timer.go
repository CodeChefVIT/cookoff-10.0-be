package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/labstack/echo/v4"
)

func GetTime(c echo.Context) error {
	roundID, err := utils.RedisClient.Get(c.Request().Context(), "current_round").Result()
	if err != nil {
		logger.Errorf("could not get current_round: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	timeStr, err := utils.RedisClient.HGet(c.Request().Context(), "round_end_time", roundID).Result()
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
		"server_time":    time.Now().Format(time.RFC3339),
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

	err = utils.RedisClient.HSet(c.Request().Context(), "round_end_time", req.RoundID, req.Time).Err()
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

	timeStr, err := utils.RedisClient.HGet(c.Request().Context(), "round_end_time", req.RoundID).Result()
	if err != nil {
		logger.Errorf("Round get time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	currTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "incorrect time format",
		})
	}

	updatedTime := currTime.Add(duration)

	err = utils.RedisClient.HSet(c.Request().Context(), "round_end_time", req.RoundID, updatedTime).Err()
	if err != nil {
		logger.Errorf("Round get time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	err = utils.RedisClient.ExpireAt(c.Request().Context(), "current_round", updatedTime).Err()
	if err != nil {
		logger.Errorf("current_round TTL updation error: %v", err.Error())
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

func StartRound(c echo.Context) error {
	round, err := utils.RedisClient.Incr(c.Request().Context(), "current_round").Result()
	if err != nil {
		logger.Errorf("could not increment current_round: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	roundID := strconv.FormatInt(round, 10)

	timeStr, err := utils.RedisClient.HGet(c.Request().Context(), "round_end_time", roundID).Result()
	if err != nil {
		logger.Errorf("Round get time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	expTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "incorrect time format",
		})
	}

	err = utils.RedisClient.Set(c.Request().Context(), "is_round_started", true, 0).Err()
	if err != nil {
		logger.Errorf("could not get current_round: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	err = utils.RedisClient.ExpireAt(c.Request().Context(), "current_round", expTime).Err()
	if err != nil {
		logger.Errorf("current_round TTL set error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":        "success",
		"message":       "round started",
		"started_round": round,
	})
}
