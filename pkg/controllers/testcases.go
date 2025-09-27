package controllers

import (
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func CreateTestCase(c echo.Context) error {
	var req dto.CreateTestCaseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid request body",
			"error":  err.Error(),
		})
	}

	if err := validator.ValidatePayload(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Validation failed",
			"error":  err.Error(),
		})
	}

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid question ID",
			"error":  err.Error(),
		})
	}

	
	var memory pgtype.Numeric
	if err := memory.Scan(req.Memory); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid memory value",
			"error":  err.Error(),
		})
	}

	var runtime pgtype.Numeric
	if err := runtime.Scan(req.Runtime); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid runtime value",
			"error":  err.Error(),
		})
	}

	testCase, err := utils.Queries.CreateTestCase(c.Request().Context(), db.CreateTestCaseParams{
		ID:             uuid.New(),
		ExpectedOutput: req.ExpectedOutput,
		Memory:         req.Memory,
		Input:          req.Input,
		Hidden:         req.Hidden,
		Runtime:        req.Runtime,
		QuestionID:     questionID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to create test case",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"status":    "success",
		"test_case": testCase,
	})
}

func GetTestCase(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid test case ID",
			"error":  err.Error(),
		})
	}

	testCase, err := utils.Queries.GetTestCase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status": "Test case not found",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":    "success",
		"test_case": testCase,
	})
}

func UpdateTestCase(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid test case ID",
			"error":  err.Error(),
		})
	}

	existing, err := utils.Queries.GetTestCase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status": "Test case not found",
			"error":  err.Error(),
		})
	}

	var req dto.UpdateTestCaseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid request body",
			"error":  err.Error(),
		})
	}

	expectedOutput := existing.ExpectedOutput
	if req.ExpectedOutput != "" {
		expectedOutput = req.ExpectedOutput
	}

	input := existing.Input
	if req.Input != "" {
		input = req.Input
	}

	hidden := existing.Hidden
	if req.Hidden != nil {
		hidden = *req.Hidden
	}

	questionID := existing.QuestionID
	if req.QuestionID != "" {
		parsedID, err := uuid.Parse(req.QuestionID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status": "Invalid question ID",
				"error":  err.Error(),
			})
		}
		questionID = parsedID
	}

	memory := existing.Memory
	if req.Memory != "" {
		var numeric pgtype.Numeric
		if err := numeric.Scan(req.Memory); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status": "Invalid memory value",
				"error":  err.Error(),
			})
		}
		memory = numeric
	}

	runtime := existing.Runtime
	if req.Runtime != "" {
		var numeric pgtype.Numeric
		if err := numeric.Scan(req.Runtime); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"status": "Invalid runtime value",
				"error":  err.Error(),
			})
		}
		runtime = numeric
	}

	updated, err := utils.Queries.UpdateTestCase(c.Request().Context(), db.UpdateTestCaseParams{
		ID:             id,
		ExpectedOutput: expectedOutput,
		Memory:         memory,
		Input:          input,
		Hidden:         hidden,
		Runtime:        runtime,
		QuestionID:     questionID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to update test case",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":    "success",
		"test_case": updated,
	})
}

func DeleteTestCase(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid test case ID",
			"error":  err.Error(),
		})
	}

	_, err = utils.Queries.GetTestCase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status": "Test case not found",
			"error":  err.Error(),
		})
	}

	err = utils.Queries.DeleteTestCase(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to delete test case",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":  "success",
		"message": "Test case deleted successfully",
	})
}

func GetTestCasesByQuestion(c echo.Context) error {
	questionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid question ID",
			"error":  err.Error(),
		})
	}

	testCases, err := utils.Queries.GetTestCasesByQuestion(c.Request().Context(), questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to fetch test cases",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":      "success",
		"test_cases":  testCases,
		"total_count": len(testCases),
	})
}

func GetPublicTestCasesByQuestion(c echo.Context) error {
	questionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Invalid question ID",
			"error":  err.Error(),
		})
	}

	testCases, err := utils.Queries.GetPublicTestCasesByQuestion(c.Request().Context(), questionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to fetch test cases",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":      "success",
		"test_cases":  testCases,
		"total_count": len(testCases),
	})
}

func GetAllTestCases(c echo.Context) error {
	testCases, err := utils.Queries.GetAllTestCases(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status": "Failed to fetch test cases",
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status":      "success",
		"test_cases":  testCases,
		"total_count": len(testCases),
	})
}
