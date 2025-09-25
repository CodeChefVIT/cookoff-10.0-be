package utils

import (
	"context"

	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
)

func InitTimer() {
	exists, err := RedisClient.Exists(context.Background(), "current_round").Result()
	if err != nil {
		logger.Errorf(err.Error())
	}

	if exists == 0 {
		err := RedisClient.Set(context.Background(), "current_round", 1, 0).Err()
		if err != nil {
			logger.Errorf(err.Error())
		}
		logger.Infof("Initialized current_round")
	} else {
		logger.Infof("current_round already exist, skipping")
	}

	exists, err = RedisClient.Exists(context.Background(), "is_round_started").Result()
	if err != nil {
		logger.Errorf(err.Error())
	}

	if exists == 0 {
		err = RedisClient.Set(context.Background(), "is_round_started", false, 0).Err()
		if err != nil {
			logger.Errorf(err.Error())
		}
		logger.Infof("Initialized is_round_started")
	} else {
		logger.Infof("is_round_started already exist, skipping")
	}
}
