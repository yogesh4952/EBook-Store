package repository

import (
	"context"
	"errors"
	"fmt"

	bookModel "github.com/yogesh4952/ebookstore/internal/book/models"
	orderModel "github.com/yogesh4952/ebookstore/internal/order/models"
	paymentModel "github.com/yogesh4952/ebookstore/internal/payment/model"
	"gorm.io/gorm"
)

type IPaymentRepo interface {
	FindOrderByTransactionUUID(ctx context.Context, transactionUUID string) (*orderModel.Order, error)
	ConfirmPayment(ctx context.Context, orderID uint, transactionCode string) error
}

type paymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *paymentRepo {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) FindOrderByTransactionUUID(ctx context.Context, transactionUUID string) (*orderModel.Order, error) {
	var order orderModel.Order

	err := r.db.WithContext(ctx).Where("transaction_uuid = ?", transactionUUID).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: %s", paymentModel.ErrOrderNotFound, transactionUUID)
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ConfirmPayment marks the order PAID, stores eSewa's transaction code as the
// external payment reference, and decrements stock for every ordered book, all
// inside a single database transaction so a failure part-way leaves the order
// untouched rather than paid-but-undelivered.
func (r *paymentRepo) ConfirmPayment(ctx context.Context, orderID uint, transactionCode string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var items []orderModel.OrderItem
		if err := tx.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
			return err
		}

		// The status guard makes this write idempotent: if a concurrent
		// callback already confirmed this order, no rows match and we leave
		// the existing PAID row alone instead of re-running fulfillment.
		result := tx.Model(&orderModel.Order{}).
			Where("id = ? AND payment_status = ?", orderID, orderModel.PayementPending).
			Updates(map[string]interface{}{
				"payment_status":   orderModel.PayementPaid,
				"transaction_code": transactionCode,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("%w: order %d is not pending", paymentModel.ErrOrderNotFound, orderID)
		}

		for _, item := range items {
			err := tx.Model(&bookModel.Book{}).
				Where("id = ? AND units >= ?", item.BookId, item.Quantity).
				Update("units", gorm.Expr("units - ?", item.Quantity)).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}
