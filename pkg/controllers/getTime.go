package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func GetTime(c echo.Context) error {
	userID, ok := c.Get(utils.UserContextKey).(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	roundStarted, err := utils.RedisClient.Get(c.Request().Context(), "is_round_started").Bool()
	if err != nil && err != redis.Nil {
		logger.Errorf("could not get is_round_started: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	if !roundStarted {
		return c.JSON(http.StatusConflict, echo.Map{
			"status": "failed",
			"error":  "round not started",
		})
	}

	roundID, err := utils.RedisClient.Get(c.Request().Context(), "current_round").Result()
	if err != nil {
		logger.Errorf("could not get current_round: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	userRound, err := utils.Queries.GetUserRound(c.Request().Context(), userID)
	if err != nil {
		logger.Errorf("could not get user's round: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}
	userRoundStr := strconv.FormatInt(int64(userRound), 10)

	if userRoundStr != roundID {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status": "failed",
			"error":  "round mismatch",
		})
	}

	roundEndTimeStr, err := utils.RedisClient.HGet(c.Request().Context(), "round_end_time", roundID).Result()
	if err != nil {
		logger.Errorf("Round get time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	roundEndTime, err := time.ParseInLocation(time.RFC3339, roundEndTimeStr, utils.IST)
	if err != nil {
		logger.Errorf("RET parse error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	roundStartTimeStr, err := utils.RedisClient.Get(c.Request().Context(), "round_start_time").Result()
	if err != nil {
		logger.Errorf("Get start round time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	roundStartTime, err := time.ParseInLocation(time.RFC3339, roundStartTimeStr, utils.IST)
	if err != nil {
		logger.Errorf("RST parse error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":           "success",
		"round_end_time":   roundEndTime,
		"server_time":      time.Now().In(utils.IST).Format(time.RFC3339),
		"round_start_time": roundStartTime,
	})
}
