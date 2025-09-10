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
}
