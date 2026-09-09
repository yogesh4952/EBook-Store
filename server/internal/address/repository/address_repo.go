package repository

import (
	"context"
	"fmt"

	"github.com/yogesh4952/ebookstore/internal/address/models"
	"gorm.io/gorm"
)

type IAddressrepo interface {
	FindUserAddressById(ctx context.Context, addreessId uint) (*models.UserAddress, error)
	AddAddress(ctx context.Context, data *models.UserAddress) error
}

type addressRepo struct {
	db *gorm.DB
}

func NewAddressRepo(db *gorm.DB) *addressRepo {
	return &addressRepo{db: db}
}

func (ar *addressRepo) AddAddress(ctx context.Context, data *models.UserAddress) error {
	err := ar.db.WithContext(ctx).Create(data).Error
	if err != nil {
		return fmt.Errorf("Error adding address")
	}
	return nil
}

func (ar *addressRepo) FindUserAddressById(ctx context.Context, id uint) (*models.UserAddress, error) {
	var address models.UserAddress
	result := ar.db.WithContext(ctx).First(&address, id)
	if result.Error != nil {
		return nil, fmt.Errorf("Invalid address Id")
	}
	return &address, nil
}
