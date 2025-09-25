package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/middlewares"
	"github.com/labstack/echo/v4"
)

func testcaseRoutes(api *echo.Group) {
	testcases := api.Group("/testcase")

	testcases.GET("/:id", controllers.GetTestCase)
	testcases.GET("", controllers.GetAllTestCases)

	// Question-specific test cases
	api.GET("/question/:id/testcases", controllers.GetTestCasesByQuestion)
	api.GET("/question/:id/testcases/public", controllers.GetPublicTestCasesByQuestion)

	// Admin only test case routes
	adminTestcases := testcases.Group("")
	adminTestcases.Use(middlewares.AdminOnly)
	{
		adminTestcases.POST("", controllers.CreateTestCase)
		adminTestcases.PUT("/:id", controllers.UpdateTestCase)
		adminTestcases.DELETE("/:id", controllers.DeleteTestCase)
	}
}
