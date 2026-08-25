package repository

import (
	"context"
	"fmt"

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

func (s *sellerRepo) FindBySellerId(ctx context., userId uint) (*models.Seller, error) {

	fmt.Println(userId)
	var seller models.Seller
	err := s.db.WithContext(ctx).Where("user_id = ? ", userId).First(&seller).Error
	if err != nil {
		return nil, fmt.Errorf("Error finding seller: %w", err)
	}

	return &seller, nil
}
