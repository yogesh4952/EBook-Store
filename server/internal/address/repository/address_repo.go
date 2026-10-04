package repository

import (
	"context"
	"fmt"

	"github.com/yogesh4952/ebookstore/internal/address/models"
	"gorm.io/gorm"
)

type IAddressrepo interface {
	AddAddress(ctx context.Context, data *models.UserAddress) error
	FindUserAddressByIdAndUser(ctx context.Context, id uint, userId uint) (*models.UserAddress, error)
	FindAllByUser(ctx context.Context, userId uint) ([]models.UserAddress, error)
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

func (ar *addressRepo) FindUserAddressByIdAndUser(ctx context.Context, id uint, userId uint) (*models.UserAddress, error) {
	var address models.UserAddress
	if err := ar.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userId).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

// repo
func (ar *addressRepo) FindAllByUser(ctx context.Context, userId uint) ([]models.UserAddress, error) {
	var addresses []models.UserAddress
	if err := ar.db.WithContext(ctx).Where("user_id = ?", userId).Find(&addresses).Error; err != nil {
		return nil, err
	}
	return addresses, nil
}
