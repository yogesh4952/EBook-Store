package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/book/services"
)

func main() {
	bookFile, err := os.Open("/home/yogesh/code/EBook-Store/server/data/book.json")

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

	err = services.IBookService.BatchBookSeed(books)
}
