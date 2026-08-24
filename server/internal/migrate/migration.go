package main

import (
	"log"

	"github.com/yogesh4952/ebookstore/internal/initializers"
	"github.com/yogesh4952/ebookstore/internal/user/models"
)

func init() {
	initializers.LoadEnv()
	initializers.InitDb()
}
func main() {
	err := initializers.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Error during automigrations: %v", err)
	} else {
		log.Print("Succesfully migrated")
	}
}
