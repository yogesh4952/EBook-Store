package repository

import (
	"context"

	"github.com/yogesh4952/ebookstore/internal/sellers/models"
	"gorm.io/gorm"
)

//contains func that calls the db

type ISeller interface {
	FindBySellerId(ctx context.Context, userId uint) (*models.Seller, error)
}

type sellerRepo struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) *sellerRepo {
	return &sellerRepo{db: db}
}

func (s *sellerRepo) FindBySellerId(ctx context.Context, userId uint) (*models.Seller, error) {

	var seller models.Seller
	err := s.db.WithContext(ctx).Where("UserID = ? ", userId).First(&seller).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}
