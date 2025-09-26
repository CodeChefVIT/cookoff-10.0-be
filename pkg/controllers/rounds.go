package controllers

import (
	"net/http"
	"strconv"
	"time"

	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func StartRound(c echo.Context) error {
	roundStr, err := utils.RedisClient.Get(c.Request().Context(), "current_round").Result()
	if err != nil {
		if err == redis.Nil {
			roundStr = "0"
		} else {
			logger.Errorf("could not get current_round: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status": "failed",
				"error":  "something went wrong",
			})
		}
	}

	val, err := strconv.Atoi(roundStr)
	if err != nil {
		logger.Errorf("int conversion err:", err)
	}

	val++
	roundID := strconv.Itoa(val)

	roundEndTimeStr, err := utils.RedisClient.HGet(c.Request().Context(), "round_end_time", roundID).Result()
	if err != nil {
		logger.Errorf("Round get time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "round end time not set",
		})
	}

	expTime, err := time.ParseInLocation(time.RFC3339, roundEndTimeStr, utils.IST)
	if err != nil {
		logger.Errorf("parse time error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "incorrect time format",
		})
	}

	err = utils.RedisClient.Set(c.Request().Context(), "round_start_time", time.Now().In(utils.IST).Format(time.RFC3339), 0).Err()
	if err != nil {
		logger.Errorf("could not set round_start_time: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
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

	err = utils.RedisClient.ExpireAt(c.Request().Context(), "is_round_started", expTime).Err()
	if err != nil {
		logger.Errorf("current_round TTL set error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	round, err := utils.RedisClient.Incr(c.Request().Context(), "current_round").Result()
	if err != nil {
		logger.Errorf("could not increment current_round: %v", err)
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

func ResetRound(c echo.Context) error {
	err := utils.RedisClient.Del(c.Request().Context(), "current_round").Err()
	if err != nil {
		logger.Errorf("could not reset current_round: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	err = utils.RedisClient.Del(c.Request().Context(), "is_round_started").Err()
	if err != nil {
		logger.Errorf("could not reset is_round_started: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	err = utils.RedisClient.Del(c.Request().Context(), "round_start_time").Err()
	if err != nil {
		logger.Errorf("could not reset round_start_time: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "reset current_round, round_start_time and is_round_started successful",
	})
}
