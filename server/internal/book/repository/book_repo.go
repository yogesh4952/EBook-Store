package repository

import (
	"context"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	"gorm.io/gorm"
)

type IBookRepo interface {
	PublishBook(ctx context.Context, book *models.Book) error
	ListBooks(ctx context.Context) ([]models.Book, error)
	SeedBooks(books []*models.Book) error
}

type bookRepo struct {
	db *gorm.DB
}

func NewBookRepo(db *gorm.DB) *bookRepo {
	return &bookRepo{db: db}
}

func (b *bookRepo) PublishBook(ctx context.Context, book *models.Book) error {
	return b.db.WithContext(ctx).Create(book).Error
}

func (b *bookRepo) ListBooks(ctx context.Context) ([]models.Book, error) {
	var books []models.Book
	result := b.db.WithContext(ctx).Find(&books)

	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil

}

func (b *bookRepo) SeedBooks(books []*models.Book) error {
	return b.db.CreateInBatches(books, 100).Error
}
