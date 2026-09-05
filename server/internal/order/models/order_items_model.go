package models

import (
	"time"

	"github.com/yogesh4952/ebookstore/internal/book/models"
)

type OrderItem struct {
	OrderId   uint `json:"order_id" gorm:"primaryKey; unique; index:idx_order_book; not null"`
	BookId    uint `json:"book_id" gorm:"primaryKey;uniqe; index:idx_order_book; not null"`
	Quantity  uint `json:"quantity" binding:"required,min=1"`
	UnitPrice uint `json:"unit_price" binding:"required,min=0"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Order Order       `json:"order,omitempty" gorm:"foreignKey:OrderId;constraint:OnDelete:CASCADE"`
	Book  models.Book `json:"book" gorm:"foreignKey:BookId"`
}
