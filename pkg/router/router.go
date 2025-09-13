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
	e.PUT("/callback", controllers.CallbackUrl)
	e.POST("/logout", controllers.Logout)
	e.POST("/refreshToken", controllers.RefreshToken)
  
	e.POST("/question", controllers.CreateQuestion)
	e.GET("/question", controllers.GetQuestion)
	e.GET("/question/:id", controllers.GetAllQuestions)
	e.PUT("/question/:id", controllers.UpdateQuestion)
	e.DELETE("/question/:id", controllers.DeleteQuestion)
	e.POST("/question/:id/bounty/activate", controllers.ActivateBounty)
  e.POST("/question/:id/bounty/deactivate", controllers.DeactivateBounty)

}
