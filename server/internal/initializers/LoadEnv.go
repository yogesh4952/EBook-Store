package initializers

import (
	"github.com/joho/godotenv"
	"github.com/yogesh4952/ebookstore/pkg/logger"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file: %v", err)
	}

}
