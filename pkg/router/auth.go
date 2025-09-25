package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/labstack/echo/v4"
)

// RegisterAuthRoutes registers routes that require authenticated users
func RegisterAuthRoutes(api *echo.Group) {
	api.POST("/logout", controllers.Logout)
	api.POST("/submit", controllers.SubmitCode)
	api.POST("/runcode", controllers.RunCode)
	api.GET("/result/:submission_id", controllers.GetResult)
	api.POST("/runcustom", controllers.RunCustom)
	api.GET("/dashboard", controllers.LoadDashboard)
}
