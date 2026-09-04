package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	bookRepo "github.com/yogesh4952/ebookstore/internal/book/repository"
	"github.com/yogesh4952/ebookstore/internal/book/service"
	"github.com/yogesh4952/ebookstore/internal/initializers"
	sellerRepo "github.com/yogesh4952/ebookstore/internal/sellers/repository"
)

func main() {
	initializers.InitDb()

	bookFile, err := os.Open("/home/yst/code/EBook-Store/server/data/book.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer bookFile.Close()
	var books []*models.BookPayload

	err = json.NewDecoder(bookFile).Decode(&books)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Loaded %d books\n", len(books))

	bRepo := bookRepo.NewBookRepo(initializers.DB)
	sRepo := sellerRepo.NewSellerRepository(initializers.DB)
	svc := service.NewBookService(bRepo, sRepo)

	err = svc.BatchBookSeed(books)
	if err != nil {
		fmt.Println(err)
	}
}
