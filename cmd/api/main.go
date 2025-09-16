package main

import (
	"fmt"
	"os"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/controllers"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/queue"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/router"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/utils"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/workers"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	logger.InitLogger()
	utils.LoadConfig()
	utils.InitCache()
	utils.InitTokenCache()
	utils.InitDB()
	validator.InitValidator()

	controllers.Queries = utils.Queries

	// Initialize queue system
	redisURI := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), "6379")
	if redisURI == ":" {
		redisURI = "localhost:6379" // Default fallback
	}

	taskServer, _ := queue.InitQueue(redisURI, 2)

	// Start the queue worker in a goroutine
	go func() {
		mux := asynq.NewServeMux()
		// Register the worker function for Judge0 callback processing
		mux.HandleFunc(queue.TypeJudge0Callback, workers.ProcessJudge0CallbackTask)

		// Start the queue server
		queue.StartQueueServer(taskServer, mux)
	}()

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogStatus:     true,
		LogError:      true,
		LogValuesFunc: logger.RouteLogger,
	}))

	router.RegisterRoute(e)

	for _, r := range e.Routes() {
		fmt.Println(r.Method, r.Path)
	}

	logger.Infof("Starting HTTP server on port %s", utils.Config.Port)
	e.Start(":" + utils.Config.Port)
}
