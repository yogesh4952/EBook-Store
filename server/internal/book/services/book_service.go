package services

import (
	"context"
	"fmt"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/book/repository"
	sellerRepo "github.com/yogesh4952/ebookstore/internal/sellers/repository"
)

type IBookService interface {
	PublishBook(ctx context.Context, userId uint, data *models.BookPayload) error
	BatchBookSeed(data []*models.BookPayload) error
}

type bookService struct {
	bookRepo   repository.IBookRepo
	sellerRepo sellerRepo.ISeller
}

func NewBookService(bookRepo repository.IBookRepo, sellerRepo sellerRepo.ISeller) *bookService {
	return &bookService{
		bookRepo:   bookRepo,
		sellerRepo: sellerRepo,
	}
}

func (b *bookService) PublishBook(ctx context.Context, userId uint, data *models.BookPayload) error {

	seller, err := b.sellerRepo.FindBySellerId(ctx, userId)
	if err != nil {
		return err
	}

	book := &models.Book{
		AuthorName:   data.AuthorName,
		Title:        data.Title,
		Genre:        data.Genre,
		Category:     data.Category,
		Pages:        data.Pages,
		Publication:  data.Publication,
		Price:        data.Price,
		Units:        data.Units,
		SellerID:     &seller.ID,
		CoverPageUrl: data.CoverPageUrl,
	}

	return b.bookRepo.PublishBook(ctx, book)

}

func (b *bookService) BatchBookSeed(payloads []*models.BookPayload) error {
	if len(payloads) < 0 {
		return fmt.Errorf("Empty data! Please provide data.")
	}

	books := make([]*models.Book, 0, len(payloads))

	for _, payload := range payloads {
		book := &models.Book{
			Title:        payload.Title,
			AuthorName:   payload.AuthorName,
			Genre:        payload.Genre,
			Category:     payload.Category,
			Pages:        payload.Pages,
			Publication:  payload.Publication,
			Price:        payload.Price,
			Units:        payload.Units,
			CoverPageUrl: payload.CoverPageUrl,
			SellerID:     payload.SellerId,
		}

		books = append(books, book)

	}

	return b.bookRepo.SeedBooks(books)

}
