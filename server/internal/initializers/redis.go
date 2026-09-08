package initializers

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/yogesh4952/ebookstore/pkg/logger"
)

var Ctx = context.Background()

func InitRedis() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	logger.Success("Successfully connected with Redis!!!")
	return client, nil
}
