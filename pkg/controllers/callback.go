package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/dto"
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

const TypeProcessSubmission = "submission:process"

func CallbackUrl(c echo.Context, taskClient *asynq.Client) error {
	logger.Infof("Judge0 Callback hit\n")

	var callbackPayload dto.Judge0CallbackPayload
	if err := c.Bind(&callbackPayload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	payload, err := json.Marshal(callbackPayload)
	if err != nil {
		logger.Errorf("Failed to marshal payload: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error queueing request",
		})
	}

	task := asynq.NewTask("submission:process", payload)
	info, err := taskClient.Enqueue(task)
	if err != nil {
		logger.Errorf("Failed to enqueue task: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error queueing request",
		})
	}
	logger.Infof("Enqueued task: %+v, Queue: %s", info.ID, info.Queue)

	log.Println("Task enqueued successfully")

	return c.NoContent(http.StatusOK)
}
