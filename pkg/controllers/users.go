package controllers

import (
	"net/http"
	"sort"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	RegNo          string    `json:"regNo"`
	Role           string    `json:"role"`
	RoundQualified int32     `json:"roundQualified"`
	Name           string    `json:"name"`
	IsBanned       bool      `json:"isBanned"`
}

func GetAllUsers(c echo.Context) error {
	users, err := utils.Queries.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not get users",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"users":  users,
	})
}

func BanUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not ban user",
			"error":  "UUID GALAT HAI BHAI",
		})
	}
	if err := utils.Queries.BanUser(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not ban user",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "User banned successfully",
	})
}

func UnbanUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Could not unban user",
			"error":  "UUID GALAT HAI BHAI",
		})
	}
	if err := utils.Queries.UnbanUser(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not unban user",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "User unbanned successfully",
	})
}

func GetLeaderboard(c echo.Context) error {
	users, err := utils.Queries.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not fetch leaderboard",
			"error":  err.Error(),
		})
	}

	type LeaderUser struct {
		db.User
		ScoreInt int64 `json:"score"`
	}

	var leaderboard []LeaderUser
	for _, u := range users {
		var scoreInt int64
		if val, err := u.Score.Int64Value(); err == nil {
			scoreInt = val.Int64
		} else {
			scoreInt = 0
		}
		leaderboard = append(leaderboard, LeaderUser{
			User:     u,
			ScoreInt: scoreInt,
		})
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].ScoreInt > leaderboard[j].ScoreInt
	})

	return c.JSON(http.StatusOK, echo.Map{
		"status":      "success",
		"leaderboard": leaderboard,
	})
}

func UpgradeUserToRound(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid UUID"})
	}

	if err := utils.Queries.UpgradeUserToRound(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not upgrade user",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "User upgraded to next round",
	})
}

func GetSubmissionByUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid UUID"})
	}

	subs, err := utils.Queries.GetSubmissionByUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not fetch submissions",
			"error":  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "success", "submissions": subs})
}