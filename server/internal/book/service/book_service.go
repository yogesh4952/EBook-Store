package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/book/repository"
	sellerRepo "github.com/yogesh4952/ebookstore/internal/sellers/repository"
	"github.com/yogesh4952/ebookstore/pkg/utils"
)

type IBookService interface {
	PublishBook(ctx context.Context, userId uint, data *models.PublishBookPayload) error
	UpdateBook(ctx context.Context, userId uint, data *models.UpdateBookPayload) (models.Book, error)
	ListBooks(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error)
	BatchBookSeed(data []*models.PublishBookPayload) error
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

func (b *bookService) PublishBook(ctx context.Context, userId uint, data *models.PublishBookPayload) error {

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

func (b *bookService) UpdateBook(ctx context.Context, userId uint, data *models.UpdateBookPayload) (models.Book, error) {
	seller, err := b.sellerRepo.FindBySellerId(ctx, userId)

	if err != nil {
		return models.Book{}, err
	}
	existingBook, err := b.bookRepo.FindById(ctx, data.BookId)

	if err != nil {
		return models.Book{}, err
	}

	if existingBook.SellerID == nil || *existingBook.SellerID != seller.ID {
		return models.Book{}, errors.New("abac vioalation: you do not have permission to update this book")
	}

	existingBook.Title = *data.Title
	existingBook.AuthorName = *data.AuthorName
	existingBook.Genre = *data.Genre
	existingBook.Category = *data.Category
	existingBook.Pages = *data.Pages
	existingBook.Publication = *data.Publication
	existingBook.Price = *data.Price
	existingBook.Units = *data.Units
	existingBook.CoverPageUrl = *data.CoverPageUrl
	return b.bookRepo.UpdateBook(ctx, existingBook)

}

func (b *bookService) ListBooks(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error) {

	return b.bookRepo.ListBooks(ctx, p)
}

func (b *bookService) BatchBookSeed(payloads []*models.PublishBookPayload) error {
	if len(payloads) == 0 {
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
