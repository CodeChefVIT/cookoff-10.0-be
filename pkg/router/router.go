// pkg/router/router.go
package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/middlewares"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func RegisterRoute(e *echo.Echo, taskClient *asynq.Client) {
	// Public routes (no auth)
	RegisterPublicRoutes(e, taskClient)

	// API group with JWT authentication
	api := e.Group("")

	// Do NOT change the order of middlewares here as the userid is set in context first and then bancheck is performed
	api.Use(middlewares.VerifyJWTMiddleware, middlewares.BanCheckUser)

	// Authenticated routes
	RegisterAuthRoutes(api)

	// Domain specific routes
	questionRoutes(api)
	testcaseRoutes(api)
	adminRoutes(api)
}
