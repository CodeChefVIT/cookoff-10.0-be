package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Logout(c echo.Context) error {
	access, err := c.Cookie("access_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "access token not found",
		})
	}

	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "refresh token not found",
		})
	}

	if access != nil {
		access.Value = ""
		access.MaxAge = -1
		access.Expires = time.Now()
		access.HttpOnly = true
		access.Secure = true
		access.Path = "/"
		access.SameSite = http.SameSiteNoneMode
		c.SetCookie(access)
	}

	if refresh != nil {
		refresh.Value = ""
		refresh.MaxAge = -1
		refresh.Expires = time.Now()
		refresh.HttpOnly = true
		refresh.Secure = true
		refresh.Path = "/"
		refresh.SameSite = http.SameSiteNoneMode
		c.SetCookie(refresh)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "logged out successfully",
	})
}
