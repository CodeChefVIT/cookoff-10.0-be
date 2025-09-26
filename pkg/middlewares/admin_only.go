package middlewares

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get(utils.UserRoleKey).(string)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status": "failed",
				"error":  "something went wrong",
			})
		}
		if role != "admin" {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status": "failed",
				"error":  "not allowed to visit",
			})
		}

		return next(c)
	}
}
