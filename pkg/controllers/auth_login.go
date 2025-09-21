package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {
	var req dto.LoginRequest

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

	user, err := utils.Queries.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, echo.Map{
				"status": "failed",
				"error":  "user not found",
			})
		}

		logger.Infof("received error from database %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "some error occurred",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return c.JSON(http.StatusConflict, echo.Map{
				"status": "failed",
				"error":  "invalid password",
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "some error occurred",
		})
	}

	accessToken, err := auth.CreateAccessToken(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "failed to generate token A",
		})
	}

	refreshToken, err := auth.CreateRefreshToken(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "failed to generate token R",
		})
	}

	expiration := (time.Hour + 25*time.Minute)
	err = utils.RedisClient.Set(c.Request().Context(), user.ID.String(), refreshToken, expiration).Err()
	if err != nil {
		logger.Errorf(fmt.Sprintf("failed to set token in cache %v", err.Error()))
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "failed to set token in cache",
		})
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   60,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   5400,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	data := echo.Map{
		"username": user.Name,
		"round":    user.RoundQualified,
		"score":    user.Score,
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "Login successful",
		"data":    data,
	})
}
