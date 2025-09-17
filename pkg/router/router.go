package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterRoute(e *echo.Echo) {
	e.GET("/ping", controllers.Ping)
	e.GET("/docs", controllers.Docs)
	e.POST("/signup", controllers.Signup)
	e.POST(("/jakabutarja"), controllers.SubmitCode)
	e.POST("/login", controllers.Login)
	e.POST("/logout", controllers.Logout)
	e.POST("/refreshToken", controllers.RefreshToken)
	e.POST("/callback", controllers.CallbackUrl)

	e.POST("/question", controllers.CreateQuestion)
	e.GET("/question", controllers.GetAllQuestions)
	e.GET("/question/:id", controllers.GetQuestion)
	e.GET("/question", controllers.GetAllQuestions)
	e.GET("/question/:id", controllers.GetQuestion)
	e.PUT("/question/:id", controllers.UpdateQuestion)
	e.DELETE("/question/:id", controllers.DeleteQuestion, middlewares.AdminOnly)
	e.POST("/question/:id/bounty/activate", controllers.ActivateBounty, middlewares.AdminOnly)
	e.POST("/question/:id/bounty/deactivate", controllers.DeactivateBounty, middlewares.AdminOnly)

	// Test case routes
	e.GET("/testcase/:id", controllers.GetTestCase)
	e.GET("/question/:id/testcases", controllers.GetTestCasesByQuestion)
	e.GET("/question/:id/testcases/public", controllers.GetPublicTestCasesByQuestion)
	e.POST("/testcase", controllers.CreateTestCase)
	e.PUT("/testcase/:id", controllers.UpdateTestCase, middlewares.AdminOnly)
	e.DELETE("/testcase/:id", controllers.DeleteTestCase, middlewares.AdminOnly)
	e.GET("/testcases", controllers.GetAllTestCases)

	// Admin Routes
	e.GET("/users", controllers.GetAllUsers, middlewares.AdminOnly)
	e.POST("/users/:id/ban", controllers.BanUser, middlewares.AdminOnly)
	e.POST("/users/:id/unban", controllers.UnbanUser, middlewares.AdminOnly)
	e.GET("/leaderboard", controllers.GetLeaderboard, middlewares.AdminOnly)
	e.POST("/users/:id/upgrade", controllers.UpgradeUserToRound, middlewares.AdminOnly)
	e.GET("/users/:id/submissions", controllers.GetSubmissionByUser, middlewares.AdminOnly)

}
