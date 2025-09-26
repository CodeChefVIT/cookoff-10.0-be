package controllers

import (
	"net/http"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func SetTime(c echo.Context) error {
	roundStarted, err := utils.RedisClient.Get(c.Request().Context(), "is_round_started").Bool()
	if err != nil && err != redis.Nil {
		logger.Errorf("could not get is_round_started: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "something went wrong",
		})
	}

	if roundStarted {
		return c.JSON(http.StatusConflict, echo.Map{
			"status": "failed",
			"error":  "round in progress",
		})
	}

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

	setTime, err := time.ParseInLocation(time.RFC3339, req.Time, utils.IST)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "failed",
			"error":  "incorrect time format",
		})
	}

	err = utils.RedisClient.HSet(c.Request().Context(), "round_end_time", req.RoundID, setTime).Err()
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
