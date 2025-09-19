package middlewares

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte(utils.Config.JwtSecret)

func VerifyJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status": "Unauthorized",
				"error":  "missing access token",
			})
		}
		tokenString := cookie.Value

		claims := &auth.AccessTokenClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status": "Unauthorized",
				"error":  "invalid or expired access token",
			})
		}

		return next(c)
	}
}

// RequireAuthExcept applies VerifyJWTMiddleware to all requests except those in the skip map.
// Keys of the map are exact path strings to skip (value is ignored).
func RequireAuthExcept(skip map[string]bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			if skip[path] {
				return next(c)
			}
			return VerifyJWTMiddleware(next)(c)
		}
	}
}
