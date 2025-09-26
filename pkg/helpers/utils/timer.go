package utils

import (
	"time"

	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
)

var IST *time.Location

func InitTimer() {
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		logger.Errorf("Timer Configuration Failed")
	}

	IST = ist
	logger.Infof("Timer Configuation Successful")
}
