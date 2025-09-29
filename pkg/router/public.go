package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

// RegisterPublicRoutes sets routes that do not require authentication
func RegisterPublicRoutes(e *echo.Echo, taskClient *asynq.Client) {
	e.GET("/ping", controllers.Ping)
	e.PUT("/callback", func(c echo.Context) error {
		return controllers.CallbackUrl(c, taskClient)
	})
	e.GET("/docs", controllers.Docs)
	e.GET("/getTime", controllers.GetTime)
	e.POST("/signup", controllers.Signup)
	e.POST("/login", controllers.Login)
	e.POST("/refreshToken", controllers.RefreshToken)
}
