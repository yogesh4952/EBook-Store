package repository

import (
	"context"

	"github.com/yogesh4952/ebookstore/internal/order/models"
	"gorm.io/gorm"
)

type IOrderRepo interface {
	PlaceOrder(ctx context.Context, data *models.Order) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *orderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) PlaceOrder(ctx context.Context, data *models.Order) error {
	return r.db.WithContext(ctx).Create(data).Error

}

func (r *orderRepo) OrderItem(ctx context.Context, data []*models.OrderItem) error {
	return r.db.WithContext(ctx).CreateInBatches(data).Error
}
