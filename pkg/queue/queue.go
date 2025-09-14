package queue

import (
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/hibiken/asynq"
)

var TaskClient *asynq.Client
var TaskServer *asynq.Server

func InitQueue(redisAddr string, concurrency int) (*asynq.Server, *asynq.Client) {
	redisConn := asynq.RedisClientOpt{Addr: redisAddr}

	server := asynq.NewServer(redisConn, asynq.Config{
		Concurrency: concurrency,
	})

	client := asynq.NewClient(redisConn)

	TaskClient = client
	TaskServer = server

	return server, client
}

func StartQueueServer(server *asynq.Server, mux *asynq.ServeMux) {
	if err := server.Run(mux); err != nil {
		logger.Errorf("Failed to start Asynq worker: %v", err)
	}
}
