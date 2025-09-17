package utils

import (
    "context"
    "fmt"
    "os"

    "github.com/redis/go-redis/v9"
)

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

func StoreToken(token string, submissionID string) error {
    if err := TokenCache.Set(Ctx, token, submissionID, 0).Err(); err != nil {
        return err
    }
    return TokenCache.SAdd(Ctx, fmt.Sprintf("submission_tokens:%s", submissionID), token).Err()
}
// get sub 
func GetSubmissionIDByToken(token string) (string, error) {
    submissionID, err := TokenCache.Get(Ctx, token).Result()
    if err == redis.Nil {
        return "", fmt.Errorf("token not found")
    }
    return submissionID, err
}
//can be used to delete all tokens 
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
