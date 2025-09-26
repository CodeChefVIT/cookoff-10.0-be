package auth

import (
	"context"
	"time"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	expirationTime := time.Now().Add(1 * time.Hour)
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
	expirationTime := time.Now().Add(2 * time.Hour)
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

func VerifyRoundAccess(ctx context.Context, userID, questionID uuid.UUID) (bool, error) {
	user, err := utils.Queries.GetUserById(ctx, userID)
	if err != nil {
		return false, err
	}
	question, err := utils.Queries.GetQuestion(ctx, questionID)
	if err != nil {
		return false, err
	}
	return user.RoundQualified >= question.Round, nil
}
