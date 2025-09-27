package middlewares

import (
	"net/http"
	"strings"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			logger.Infof("invalid user id %v", claims.UserID)
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"status": "Unauthorized",
				"error":  "invalid user id",
			})
		}

		c.Set(utils.UserContextKey, userID)

		role := strings.ToLower(claims.Role)
		c.Set(utils.UserRoleKey, role)

		return next(c)
	}
}
