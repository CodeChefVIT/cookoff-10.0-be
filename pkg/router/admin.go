package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/middlewares"
	"github.com/labstack/echo/v4"
)

func adminRoutes(api *echo.Group) {
	admin := api.Group("/admin")
	admin.Use(middlewares.AdminOnly)

	// User management
	admin.GET("/users", controllers.GetAllUsers)
	admin.POST("/users/:id/ban", controllers.BanUser)
	admin.POST("/users/:id/unban", controllers.UnbanUser)
	admin.POST("/users/:id/upgrade", controllers.UpgradeUserToRound)
	admin.GET("/users/:id/submissions", controllers.GetUserSubmissions)

	// Leaderboard
	admin.GET("/leaderboard", controllers.GetLeaderboard)

	// Analytics
	admin.GET("/analytics", controllers.GetAnalytics)

	//Timer
	admin.POST("/setTime", controllers.SetTime)
	admin.POST("/updateTime", controllers.UpdateTime)
	admin.GET("/startRound", controllers.StartRound)
	admin.GET("/resetRound", controllers.ResetRound)
}
