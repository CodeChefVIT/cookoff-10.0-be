package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CallbackUrl(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Errorf("failed to read request body: %w", err).Error(),
		})
	}

	fmt.Printf("Judge0 Callback JSON: %s\n", string(body))
	return c.NoContent(http.StatusOK)
}
