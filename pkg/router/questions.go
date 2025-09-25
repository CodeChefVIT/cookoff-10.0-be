package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/middlewares"
	"github.com/labstack/echo/v4"
)

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
}
