package initializers

import (
	"fmt"
	"os"
	"time"

	"github.com/yogesh4952/ebookstore/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func InitDb() (*gorm.DB, error) {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "postgres://yogesh@localhost:5432/ebookstore"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.NewGormLogger(gormlogger.Warn, 200*time.Millisecond),
	})
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %w", err)
	}

	logger.Success("Successfully connected with DB")
	return db, nil
}
