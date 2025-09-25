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
	admin.GET("/users/:id/submissions", controllers.GetSubmissionByUser)

	// Leaderboard
	admin.GET("/leaderboard", controllers.GetLeaderboard)
}
