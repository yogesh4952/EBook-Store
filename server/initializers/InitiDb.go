package initializers

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	var err error
	dsn := "postgres://yogesh@localhost:5432/ebookstore"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting with db %v", err)
	} else {
		log.Print("Succesfully connected with db")
	}
}
