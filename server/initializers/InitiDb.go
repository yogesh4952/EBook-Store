package initializers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/yogesh4952/ebookstore/pkg/logger"
)

var DB *gorm.DB

func InitDb() {
	var err error
	dsn := "postgres://yogesh@localhost:5432/ebookstore"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Error connecting with db: %v", err)
	} else {
		logger.Success("Successfully connected with DB")
	}
}
