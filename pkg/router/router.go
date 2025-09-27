// pkg/router/router.go
package router

import (
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/middlewares"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func RegisterRoute(e *echo.Echo, taskClient *asynq.Client) {
	// Public routes (no auth)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(10))))

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
