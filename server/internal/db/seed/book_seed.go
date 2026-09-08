package main

import (
	"encoding/json"
	"os"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	bookRepo "github.com/yogesh4952/ebookstore/internal/book/repository"
	"github.com/yogesh4952/ebookstore/internal/book/service"
	"github.com/yogesh4952/ebookstore/internal/initializers"
	sellerRepo "github.com/yogesh4952/ebookstore/internal/sellers/repository"
	"github.com/yogesh4952/ebookstore/pkg/logger"
)

func main() {
	logger.Init()
	initializers.LoadEnv()

	db, err := initializers.InitDb()
	if err != nil {
		logger.Fatal("Failed to initialize DB: %v", err)
	}

	bookFile, err := os.Open("/home/yst/code/EBook-Store/server/data/book.json")

	if err != nil {
		logger.Error("%v", err)
		return
	}

	defer bookFile.Close()
	var books []*models.PublishBookPayload

	err = json.NewDecoder(bookFile).Decode(&books)

	if err != nil {
		logger.Error("%v", err)
		return
	}

	logger.Info("Loaded %d books", len(books))

	bRepo := bookRepo.NewBookRepo(db)
	sRepo := sellerRepo.NewSellerRepository(db)
	svc := service.NewBookService(bRepo, sRepo)

	err = svc.BatchBookSeed(books)
	if err != nil {
		logger.Error("Failed to seed books: %v", err)
	}
}
