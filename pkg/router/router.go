package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func RegisterRoute(e *echo.Echo) {
	e.GET("/ping", controllers.Ping)
	e.GET("/docs", controllers.Docs)
	e.POST("/signup", controllers.Signup)
	e.POST(("/jakabutarja"), controllers.SubmitCode)
	e.POST("/login", controllers.Login)

	e.POST("/question", controllers.CreateQuestion)
	e.GET("/question", controllers.GetAllQuestions)
	e.GET("/question/:id", controllers.GetQuestion)
	e.PUT("/question/:id", controllers.UpdateQuestion)
	e.DELETE("/question/:id", controllers.DeleteQuestion)
	e.POST("/question/:id/bounty/activate", controllers.ActivateBounty)
	e.POST("/question/:id/bounty/deactivate", controllers.DeactivateBounty)

	e.GET("/users", controllers.GetAllUsers)
	e.POST("/users/:id/ban", controllers.BanUser)
	e.POST("/users/:id/unban", controllers.UnbanUser)
	e.GET("/leaderboard", controllers.GetLeaderboard)
	e.POST("/users/:id/upgrade", controllers.UpgradeUserToRound)
	e.GET("/users/:id/submissions", controllers.GetSubmissionByUser)

}
