package controllers

import (
	"math/big"
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/auth"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c echo.Context) error {
	var payload dto.SignupRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "failed",
			"error":  "invalid request",
		})
	}

	if err := validator.ValidatePayload(payload); err != nil {
		return c.JSON(http.StatusNotAcceptable, echo.Map{
			"status": "failed",
			"error":  "invalid input",
		})
	}

	if payload.Key != utils.Config.SecretKey {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"status": "failed",
			"error":  "stfu",
		})
	}

	_, err := utils.Queries.GetUserByEmail(c.Request().Context(), payload.Email)
	if err == nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"status": "failed",
			"error":  "User already exists",
		})
	}

	password := auth.PasswordGenerator(6)
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "some error occurred while generation",
		})
	}

	id, err := uuid.NewV7()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "some error occurred wihle uuid generation",
		})
	}

	_, err = utils.Queries.CreateUser(c.Request().Context(), db.CreateUserParams{
		ID:             id,
		Email:          payload.Email,
		RegNo:          payload.RegNo,
		Password:       string(hashed),
		Role:           "user",
		RoundQualified: 0,
		Score: pgtype.Numeric{
			Int:              big.NewInt(0),
			Exp:              0,
			NaN:              false,
			InfinityModifier: 0,
			Valid:            true,
		},
		Name: payload.Name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "failed",
			"error":  "some error occurred while creating user",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":   "success",
		"message":  "user added",
		"email":    payload.Email,
		"password": password,
	})
}
