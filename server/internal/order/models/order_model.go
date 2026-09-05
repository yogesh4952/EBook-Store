package models

import (
	"github.com/yogesh4952/ebookstore/internal/user/models"
	"gorm.io/gorm"
)

type PaymentMethod string
type PaymentStatus string
type OrderStatus string

const (
	PayementCOD   PaymentMethod = "COD"
	PayementEsewa PaymentMethod = "ESEWA"
)

const (
	PayementPaid    PaymentStatus = "PAID"
	PayementPending PaymentStatus = "PENDING"
	PayementRefund  PaymentStatus = "REFUND"
)
const (
	OrderDelivered    OrderStatus = "DELIVERED"
	PayementCancelled OrderStatus = "CANCELLED"
	OrderPending      OrderStatus = "PENDING"
)

type Order struct {
	gorm.Model
	OrderCode     string        `json:"order_code" gorm:"unique;not null"`
	PaymentMethod PaymentMethod `json:"payment_method" gorm:"check:payment_method IN ('COD','ESEWA');not null;"`
	PaymentStatus PaymentStatus `json:"payment_status" gorm:"check:payment_status IN ('PAID','PENDING','REFUND');not null; default:'PENDING'"`
	OrderStatus   OrderStatus   `json:"order_status" gorm:"check:order_status IN ('PAID','DELIVERED','CANCELLED' ,'PENDING');not null; default:'PENDING'"`
	TotalPrice    uint          `json:"total_price" binding:"required,min=1"`
	UserId        uint          `json:"user_id"`
	User          models.User   `gorm:"foreignKey:UserId"`

	OrderItem []OrderItem `json:"order_items"`
}
