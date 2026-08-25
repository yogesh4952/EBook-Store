package main

import (
	"log"

	bookModels "github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/initializers"
	sellerModels "github.com/yogesh4952/ebookstore/internal/sellers/models"
	userModels "github.com/yogesh4952/ebookstore/internal/user/models"
)

func init() {
	initializers.LoadEnv()
	initializers.InitDb()
}
func main() {
	err := initializers.DB.AutoMigrate(&userModels.User{}, &sellerModels.Seller{}, &bookModels.Book{})
	if err != nil {
		log.Fatalf("Error during automigrations: %v", err)
	} else {
		log.Print("Succesfully migrated")
	}
}
