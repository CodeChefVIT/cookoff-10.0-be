package workers

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
)

// ProcessJudge0CallbackTask simply logs the received Judge0 callback
func ProcessJudge0CallbackTask(ctx context.Context, t *asynq.Task) error {
	log.Printf("Received Judge0 callback task: %s", string(t.Payload()))
	return nil
}
