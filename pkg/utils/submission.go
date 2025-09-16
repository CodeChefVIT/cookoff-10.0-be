package utils

import (
	"context"
	"fmt"
	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

)

type SubmissionInput struct {
	ID         uuid.UUID
	QuestionID string
	LanguageID int
	SourceCode string
	UserID     string
}
// place holder code -- renove when queue is setup

func SaveSubmission(sub SubmissionInput) error {
	if Queries == nil {
		return fmt.Errorf("DB Queries not initialized")
	}

	qid, err := uuid.Parse(sub.QuestionID)
	if err != nil {
		return fmt.Errorf("invalid QuestionID: %v", err)
	}

	userID, err := uuid.Parse(sub.UserID)
	if err != nil {
		return fmt.Errorf("invalid UserID: %v", err)
	}

	status := "pending"

	ctx := context.Background()
	err = Queries.CreateSubmission(ctx, db.CreateSubmissionParams{
		ID:         sub.ID,
		QuestionID: qid,
		LanguageID: int32(sub.LanguageID),
		SourceCode: sub.SourceCode,
		Status:     &status,
		UserID:     uuid.NullUUID{UUID: userID, Valid: true},
		
	})
	return err
}
func GetSubmissionIDByToken(token string) (string, error) {
	if RedisClient == nil {
		return "", errors.New("redis client not initialized")
	}

	ctx := context.Background()
	submissionID, err := RedisClient.Get(ctx, "token:"+token).Result()
	if err == redis.Nil {
		return "", errors.New("token not found")
	} else if err != nil {
		return "", err
	}
	return submissionID, nil
}