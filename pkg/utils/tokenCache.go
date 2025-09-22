package utils

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/CodeChefVIT/cookoff-10.0-be/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

const SubmissionDoneStatus = "DONE"

var (
	TokenCache *redis.Client
	Ctx        = context.Background()
)

func InitTokenCache() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	TokenCache = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   1,
	})

	if err := TokenCache.Ping(Ctx).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to Redis token cache: %v", err))
	}
}

func StoreToken(token string, submissionID string, testcaseID string) error {
	if err := TokenCache.Set(Ctx, token, fmt.Sprintf("%s:%s", submissionID, testcaseID), 0).Err(); err != nil {
		return err
	}
	return TokenCache.SAdd(Ctx, fmt.Sprintf("submission_tokens:%s", submissionID), token).Err()
}

// get sub
func GetSubmissionIDByToken(ctx context.Context, token string) (string, string, error) {
	subId, err := TokenCache.Get(ctx, fmt.Sprintf("token:%s", token)).Result()
	if err == redis.Nil {
		return "", "", fmt.Errorf("token not found")
	} else if err != nil {
		return "", "", err
	}
	temp := strings.Split(subId, ":")
	return temp[0], temp[1], nil
}

// can be used to delete all tokens
func DeleteTokensBySubmissionID(submissionID string) error {
	setKey := fmt.Sprintf("submission_tokens:%s", submissionID)
	tokens, err := TokenCache.SMembers(Ctx, setKey).Result()
	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		if err := TokenCache.Del(Ctx, tokens...).Err(); err != nil {
			return err
		}
	}

	return TokenCache.Del(Ctx, setKey).Err()
}

func DeleteToken(ctx context.Context, token string) error {
	subID, _, err := GetSubmissionIDByToken(ctx, token)
	if err != nil {
		return err
	}

	err = TokenCache.Del(ctx, fmt.Sprintf("token:%s", token)).Err()
	if err != nil {
		return err
	}

	err = TokenCache.SRem(ctx, fmt.Sprintf("sub:%s:tokens", subID), token).Err()
	if err != nil {
		return err
	}

	setSize, err := TokenCache.SCard(ctx, fmt.Sprintf("sub:%s:tokens", subID)).Result()
	if err != nil {
		return err
	}

	if setSize == 0 {
		err = TokenCache.Del(ctx, fmt.Sprintf("sub:%s:tokens", subID)).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetTokenCount(ctx context.Context, subID string) (int64, error) {
	tokenCount, err := TokenCache.SCard(ctx, fmt.Sprintf("sub:%s:tokens", subID)).Result()
	if err != nil {
		return 0, err
	}
	return tokenCount, nil
}

func UpdateSubmission(ctx context.Context, idUUID uuid.UUID) error {
	status := SubmissionDoneStatus

	data, err := Queries.GetStatsForFinalSubEntry(ctx, idUUID)
	if err != nil {
		log.Println("Error Fetching submission results: ", err)
		return err
	}
	var runtime float64
	var memory int64
	var passed, failed int
	for _, v := range data {
		temp, err := v.Runtime.Float64Value()
		if err != nil {
			log.Println("Failed to convert runtime to float kms")
			return err
		}
		runtime += temp.Float64
		memory += v.Memory.Int.Int64()
		if v.Status == "success" {
			passed += 1
		} else {
			failed += 1
		}
	}

	err = Queries.UpdateSubmission(ctx, db.UpdateSubmissionParams{
		Runtime:         pgtype.Numeric{Int: big.NewInt(int64(runtime)), Valid: true},
		Memory:          pgtype.Numeric{Int: big.NewInt(int64(memory)), Valid: true},
		Status:          &status,
		TestcasesPassed: pgtype.Int4{Int32: int32(passed), Valid: true},
		TestcasesFailed: pgtype.Int4{Int32: int32(failed), Valid: true},
		ID:              idUUID,
	})

	if err != nil {
		log.Println("Error updating submission: ", err)
		return err
	}

	err = Queries.UpdateScore(ctx, idUUID)
	if err != nil {
		log.Println("Error updating the score: ", err)
		return err
	}

	log.Printf("Submission ID: %v\n", idUUID)
	return nil
}
