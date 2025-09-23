// pkg/router/router.go
package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/middlewares"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func RegisterRoute(e *echo.Echo, taskClient *asynq.Client) {
	e.GET("/ping", controllers.Ping)
	e.PUT("/callback", func(c echo.Context) error {
		return controllers.CallbackUrl(c, taskClient)
	})
	e.GET("/docs", controllers.Docs)
	e.POST("/signup", controllers.Signup)
	e.POST("/login", controllers.Login)
	e.POST("/refreshToken", controllers.RefreshToken)
	

	// e.POST("/submit", controllers.SubmitCode)

	// API group with JWT authentication
	api := e.Group("")

	// Do NOT change the order of middlewares here as the userid is set in context first and then bancheck is performed
	api.Use(middlewares.VerifyJWTMiddleware, middlewares.BanCheckUser)

	// Authenticated user routes
	api.POST("/logout", controllers.Logout)
	api.POST("/submit", controllers.SubmitCode)
	api.POST("/runcode", controllers.RunCode)
	api.GET("/result/:submission_id", controllers.GetResult)
	api.POST("/runcustom", controllers.RunCustom)
	api.GET("/dashboard", controllers.LoadDashboard)

	// Question routes
	questionRoutes(api)

	// Admin routes
	adminRoutes(api)
}

func questionRoutes(api *echo.Group) {
	questions := api.Group("/question")

	questions.GET("/round", controllers.GetQuestionsByRound)

	// Admin only question routes
	adminQuestions := questions.Group("")

	// Do NOT change the order of middlewares here as the userid is set in context first and then bancheck is performed
	adminQuestions.Use(middlewares.AdminOnly, middlewares.BanCheckUser)
	{
		adminQuestions.GET("", controllers.GetAllQuestions)
		adminQuestions.GET("/:id", controllers.GetQuestion)
		adminQuestions.POST("", controllers.CreateQuestion)
		adminQuestions.PUT("/:id", controllers.UpdateQuestion)
		adminQuestions.DELETE("/:id", controllers.DeleteQuestion)
		adminQuestions.POST("/:id/bounty/activate", controllers.ActivateBounty)
		adminQuestions.POST("/:id/bounty/deactivate", controllers.DeactivateBounty)
	}

	// Test case routes
	testcaseRoutes(api)
}

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

func adminRoutes(api *echo.Group) {
	admin := api.Group("/admin")
	admin.Use(middlewares.AdminOnly)

	// User management
	admin.GET("/users", controllers.GetAllUsers)
	admin.POST("/users/:id/ban", controllers.BanUser)
	admin.POST("/users/:id/unban", controllers.UnbanUser)
	admin.POST("/users/:id/upgrade", controllers.UpgradeUserToRound)
	admin.GET("/users/:id/submissions", controllers.GetSubmissionByUser)

	// Leaderboard
	admin.GET("/leaderboard", controllers.GetLeaderboard)
}
