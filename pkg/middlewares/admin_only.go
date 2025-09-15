package middlewares

import (
	"net/http"
	"strings"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status": "failed",
				"error":  "missing access token",
			})
		}

		claims := &auth.AccessTokenClaims{}

		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status": "failed",
				"error":  "invalid or expired token",
			})
		}

		if strings.ToLower(claims.Role) != "admin" {
			return c.JSON(http.StatusForbidden, echo.Map{
				"status": "failed",
				"error":  "not allowed to visit",
			})
		}

		return next(c)
	}
}
