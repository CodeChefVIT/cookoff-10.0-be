package middlewares

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func BanCheckUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, ok := c.Get(utils.UserContextKey).(uuid.UUID)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status": "Unauthorized",
				"error":  "user id is not present in context",
			})
		}

		user, err := utils.Queries.GetUserById(c.Request().Context(), userID)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status": "Unauthorized",
				"error":  "User was not found in database",
			})
		}

		if user.IsBanned {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status": "Unauthorized",
				"error":  "User is banned",
			})
		}

		return next(c)
	}
}
