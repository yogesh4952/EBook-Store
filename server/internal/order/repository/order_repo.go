package repository

import (
	"context"

	"github.com/yogesh4952/ebookstore/internal/order/models"
	"gorm.io/gorm"
)

type IOrderRepo interface {
	PlaceOrder(ctx context.Context, data *models.Order, items []*models.OrderItem) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *orderRepo {
	return &orderRepo{db: db}
}

const orderItemBatchSize = 100

func (r *orderRepo) PlaceOrder(ctx context.Context, data *models.Order, items []*models.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(data).Error; err != nil {
			return err
		}

		if len(items) == 0 {
			return nil
		}

		for _, item := range items {
			item.OrderId = data.ID
		}

		return tx.CreateInBatches(items, orderItemBatchSize).Error
	})
}
