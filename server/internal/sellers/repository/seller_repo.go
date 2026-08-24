package repository

import (
	"context"
	"fmt"

	"github.com/yogesh4952/ebookstore/internal/sellers/models"
	"gorm.io/gorm"
)

//contains func that calls the db

type ISeller interface {
}

type sellerRepo struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) *sellerRepo {
	return &sellerRepo{db: db}
}

func (s *sellerRepo) PublishBook(ctx context.Context, seller *models.Seller) error {

	err := s.db.WithContext(ctx).Create(seller).Error
	if err != nil {
		return err
	}

	return fmt.Errorf("Failed to create seller in db: %w", err)
}
