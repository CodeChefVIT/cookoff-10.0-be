package main

import (
	"fmt"
	"os"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/utils"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/helpers/validator"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/queue"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/workers"
	"github.com/hibiken/asynq"
)

func main() {
	logger.InitLogger()
	utils.LoadConfig()
	utils.InitCache()
	utils.InitTokenCache()
	utils.InitDB()
	validator.InitValidator()
	utils.InitTimer()

	redisURI := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), "6379")
	if redisURI == ":" {
		redisURI = "localhost:6379"
	}

	taskServer, _ := queue.InitQueue(redisURI, utils.Config.WorkerCount)

	mux := asynq.NewServeMux()
	mux.HandleFunc("submission:process", workers.ProcessJudge0CallbackTask)
	queue.StartQueueServer(taskServer, mux)
}
