package controllers

import (
	"net/http"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Logout(c echo.Context) error {
	access, err := c.Cookie("access_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "failed",
			"error":  "access token not found",
		})
	}

	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "failed",
			"error":  "refresh token not found",
		})
	}

	if refresh != nil {
		refreshToken := refresh.Value
		claims := &auth.RefreshTokenClaims{}
		token, _ := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
			return utils.Config.JwtSecret, nil
		})
		if token != nil && token.Valid {
			err := utils.RedisClient.Del(c.Request().Context(), claims.UserID).Err()
			if err != nil {
				logger.Errorf("Redis token del err: %v", err.Error())
			}
		}
	}

	if access != nil {
		access.Value = ""
		access.MaxAge = -1
		access.Expires = time.Now()
		access.HttpOnly = true
		access.Secure = utils.Config.CookieSecure
		access.Path = "/"
		access.SameSite = http.SameSiteNoneMode
		c.SetCookie(access)
	}

	if refresh != nil {
		refresh.Value = ""
		refresh.MaxAge = -1
		refresh.Expires = time.Now()
		refresh.HttpOnly = true
		refresh.Secure = utils.Config.CookieSecure
		refresh.Path = "/"
		refresh.SameSite = http.SameSiteNoneMode
		c.SetCookie(refresh)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "logged out successfully",
	})
}
