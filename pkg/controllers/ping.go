package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	fmt.Println(c.RealIP())
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Pong",
	})
}
