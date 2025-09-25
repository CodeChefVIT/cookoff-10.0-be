package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

var jwtSecret = []byte(utils.Config.JwtSecret)

func RefreshToken(c echo.Context) error {
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status": "failed",
			"error":  "refresh token missing",
		})
	}
	refreshToken := refreshCookie.Value

	token, err := jwt.ParseWithClaims(refreshToken, &auth.RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		logger.Errorf("Refrsh Token Parsing error: %v", err.Error())
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status": "failed",
			"error":  "invalid refresh token",
		})
	}

	claims, ok := token.Claims.(*auth.RefreshTokenClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status": "failed",
			"error":  "invalid claims",
		})
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status": "failed",
			"error":  "invalid userID",
		})
	}

	storedToken, err := utils.RedisClient.Get(c.Request().Context(), claims.UserID).Result()
	if errors.Is(err, redis.Nil) {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status": "failed",
			"error":  "refresh token expired",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "some error occurred",
		})
	}

	if storedToken != refreshToken {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status": "failed",
			"error":  "refresh token mismatch",
		})
	}

	user, err := utils.Queries.GetUserById(c.Request().Context(), userID)
	if err != nil {
		logger.Errorf("DB error: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "some error occurred",
		})
	}

	accessToken, err := auth.CreateAccessToken(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "failed to generate access token",
		})
	}

	newRefreshToken, err := auth.CreateRefreshToken(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "failed to generate refresh token",
		})
	}

	expiration := (time.Hour + 25*time.Minute)
	err = utils.RedisClient.Set(c.Request().Context(), claims.UserID, newRefreshToken, expiration).Err()
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
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		MaxAge:   7200,
		HttpOnly: true,
		Secure:   utils.Config.CookieSecure,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "tokens refreshed",
	})
}
