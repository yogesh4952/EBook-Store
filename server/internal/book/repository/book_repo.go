package repository

import (
	"context"

	"github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/pkg/utils"
	"gorm.io/gorm"
)

type IBookRepo interface {
	PublishBook(ctx context.Context, book *models.Book) error
	ListBooks(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error)
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
func Paginate(p utils.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.Getlimit())
	}
}

func (b *bookRepo) ListBooks(ctx context.Context, p utils.Pagination) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	if err := b.db.WithContext(ctx).Model(&models.Book{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := b.db.WithContext(ctx).
		Scopes(Paginate(p)).
		Preload("Seller").
		Preload("Seller.User").
		Find(&books).Error

	if err != nil {
		return nil, 0, err
	}

	return books, total, nil

}

func (b *bookRepo) SeedBooks(books []*models.Book) error {
	return b.db.CreateInBatches(books, 100).Error
}
