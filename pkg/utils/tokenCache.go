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

// InitTokenCache initializes a separate Redis client for token cache
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
        DB:   1, // Use DB 1 to keep it separate from queue DB 0
    })

    // Test connection
    if err := TokenCache.Ping(Ctx).Err(); err != nil {
        panic(fmt.Sprintf("failed to connect to Redis token cache: %v", err))
    }
}
