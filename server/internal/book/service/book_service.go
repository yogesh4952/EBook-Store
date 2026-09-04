package service

import (
	"context"
	"fmt"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/book/repository"
	sellerRepo "github.com/yogesh4952/ebookstore/internal/sellers/repository"
	"github.com/yogesh4952/ebookstore/pkg/logger"
	"github.com/yogesh4952/ebookstore/pkg/utils"
)

type IBookService interface {
	PublishBook(ctx context.Context, userId uint, data *models.PublishBookPayload) error
	UpdateBook(ctx context.Context, userId uint, data *models.UpdateBookPayload) (*models.Book, error)
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

func (b *bookService) UpdateBook(ctx context.Context, userId uint, data *models.UpdateBookPayload) (*models.Book, error) {

	if data == nil || data.IsEmpty() {

		return nil, fmt.Errorf("Empty payload")

	}
	updates := make(map[string]interface{})
	if data.Title != nil {
		updates["title"] = *data.Title
	}
	if data.AuthorName != nil {
		updates["author_name"] = *data.AuthorName
	}
	if data.Genre != nil {
		updates["genre"] = *data.Genre
	}
	if data.Category != nil {
		updates["category"] = *data.Category
	}
	if data.Pages != nil {
		updates["pages"] = *data.Pages
	}
	if data.Publication != nil {
		updates["publication"] = *data.Publication
	}
	if data.Units != nil {
		updates["units"] = *data.Units
	}
	if data.Price != nil {
		updates["price"] = *data.Price
	}
	if data.CoverPageUrl != nil {
		updates["cover_page_url"] = *data.CoverPageUrl
	}

	updatedBook, err := b.bookRepo.UpdateBookAtomic(ctx, userId, data.BookId, updates)
	if err != nil {
		logger.Ctx(ctx).Error().
			Uint("user_id", userId).
			Uint("book_id", data.BookId).
			Err(err).
			Msg("book update failed")
		return &models.Book{}, err
	}

	logger.Ctx(ctx).Info().
		Uint("user_id", userId).
		Uint("book_id", data.BookId).
		Msg("book update completed")

	return updatedBook, nil

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
