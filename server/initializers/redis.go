package initializers

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/yogesh4952/ebookstore/pkg/logger"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		logger.Error("Failed to connect with redis: %v", err)
	}

	logger.Success("Successfully connected with Redis!!!")
}
