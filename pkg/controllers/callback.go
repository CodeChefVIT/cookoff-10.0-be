package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/queue"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func CallbackUrl(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Errorf("failed to read request body: %w", err).Error(),
		})
	}

	fmt.Printf("Judge0 Callback JSON: %s\n", string(body))

	// Parse the Judge0 callback payload
	var callbackPayload queue.Judge0CallbackPayload
	if err := json.Unmarshal(body, &callbackPayload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Errorf("failed to parse Judge0 callback payload: %w", err).Error(),
		})
	}

	// Enqueue the callback for processing
	if queue.TaskClient != nil {
		// Serialize payload to JSON bytes
		payloadBytes, err := json.Marshal(callbackPayload)
		if err != nil {
			fmt.Printf("Failed to marshal Judge0 callback payload: %v\n", err)
		} else {
			task := asynq.NewTask(queue.TypeJudge0Callback, payloadBytes)
			_, err = queue.TaskClient.Enqueue(task, asynq.Queue("judge0"))
			if err != nil {
				fmt.Printf("Failed to enqueue Judge0 callback: %v\n", err)
				// Don't return error to Judge0, just log it
			} else {
				fmt.Printf("Successfully enqueued Judge0 callback for submission %s\n", callbackPayload.SubmissionID)
			}
		}
	} else {
		fmt.Println("TaskClient not initialized, skipping enqueue")
	}

	return c.NoContent(http.StatusOK)
}
