package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	"github.com/labstack/echo/v4"
)

// Hardcoded language names for Judge0 IDs
var languageNames = map[int32]string{
	50: "C",
	54: "C++",
	60: "Go (🥰)",
	73: "Rust",
	63: "Javascript🤮",
	51: "C#",
	62: "Java",
	68: "PHP",
	71: "Python",
}

// Reverse map: language name → ID
var languageIDs = make(map[string]int32)

func init() {
	for id, name := range languageNames {
		languageIDs[name] = id
	}
}

// GetAnalytics returns analytics data like total users, submissions, round-wise data, and language stats
func GetAnalytics(c echo.Context) error {
	ctx := c.Request().Context()

	// Total users
	userCount, err := utils.Queries.GetTotalUsersCount(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not fetch analytics",
			"error":  err.Error(),
		})
	}

	// Total submissions
	submissionCount, err := utils.Queries.GetTotalSubmissionsCount(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not fetch submissions count",
			"error":  err.Error(),
		})
	}

	// Round-wise question submissions
	roundWise, err := utils.Queries.GetRoundWiseQuestionSubmissions(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not fetch round-wise submissions",
			"error":  err.Error(),
		})
	}

	roundWiseData := make(map[string][]map[string]interface{})
	for _, r := range roundWise {
		rid := strconv.Itoa(int(r.RoundID)) // convert int32 → int → string
		roundWiseData[rid] = append(roundWiseData[rid], map[string]interface{}{
			"question_id":      r.QuestionID,
			"submissions_made": r.SubmissionsCount,
		})
	}

	// Language-wise submissions
	langs, err := utils.Queries.GetSubmissionsByLanguage(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Could not fetch language submissions",
			"error":  err.Error(),
		})
	}

	langData := make(map[string]int64)
	for _, l := range langs {
		name, ok := languageNames[l.LanguageID]
		if !ok {
			// fallback for unknown Judge0 IDs
			name = fmt.Sprintf("Unknown (%d)", l.LanguageID)
		}
		langData[name] = l.SubmissionsCount
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":            "success",
		"total_users":       userCount,
		"total_submissions": submissionCount,
		"round_wise":        roundWiseData,
		"language_wise":     langData,
	})
}
