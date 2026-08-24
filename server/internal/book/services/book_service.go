package services

import (
	"context"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/book/repository"
	sellerRepo "github.com/yogesh4952/ebookstore/internal/sellers/repository"
)

type IBookService interface {
	PublishBook(ctx context.Context, userId uint, data *models.BookPayload) error
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
	// query weather seller exist or not

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
