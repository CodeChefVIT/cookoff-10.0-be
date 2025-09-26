package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func timerRoutes(api *echo.Group) {
	time := api.Group("timer")

	time.GET("/getTime/", controllers.GetTime)
	time.POST("/setTime", controllers.SetTime)
	time.POST("/updateTime", controllers.UpdateTime)
}
