package auth

import (
	"fmt"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte(utils.Config.JwtSecret)

type AccessTokenClaims struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func CreateAccessToken(user *db.User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Minute)
	claims := &AccessTokenClaims{
		Username: user.Name,
		UserID:   user.ID.String(),
		Role:     user.Role,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func CreateRefreshToken(user *db.User) (string, error) {
	expirationTime := time.Now().Add(1*time.Hour + 30*time.Minute)
	claims := &RefreshTokenClaims{
		UserID: user.ID.String(),
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GetUserID(c echo.Context) (uuid.UUID, error) {
	user := c.Get("user").(*jwt.Token)

	if user == nil {
		return uuid.UUID{}, fmt.Errorf("JWT Token not found")
	}

	claims := user.Claims.(jwt.MapClaims)

	if !user.Valid {
		return uuid.UUID{}, fmt.Errorf("the JWT Token is invalid")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("user_id not found in token")
	}

	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user_id format: %v", err)
	}

	return uid, nil
}
