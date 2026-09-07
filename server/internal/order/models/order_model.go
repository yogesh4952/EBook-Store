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
	UserId        uint          `json:"user_id" gorm:"index:idx_userid"`

	ShippingCity            string `json:"shipping_city" binding:"required" gorm:"not null"`
	ShippingDeliveryAddress string `json:"shipping_delivery_address" binding:"required" gorm:"not null"`

	User      models.User `gorm:"foreignKey:UserId"`
	OrderItem []OrderItem `json:"order_items"`
}

type OrderItemRequest struct {
	BookId   uint `json:"book_id" binding:"required"`
	Quantity uint `json:"quantity" binding:"required,min=1"`
	SellerId uint `json:"seller_id" binding:"required"`
}

type PlaceOrderPayload struct {
	PaymentMethod PaymentMethod      `json:"payment_method" binding:"required"`
	UserAddressId uint               `json:"user_address_id" binding:"required"`
	Items         []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type OrderItemResponse struct {
	BookId    uint    `json:"book_id"`
	Title     string  `json:"title"`
	Quantity  uint    `json:"quantity"`
	UnitPrice float32 `json:"unit_price"`
	Subtotal  float32 `json:"subtotal"`
}
type PlaceOrderResponse struct {
	Success       bool                `json:"success"`
	Message       string              `json:"message"`
	OrderID       uint                `json:"order_id"`
	OrderCode     string              `json:"order_code"` // Human-readable order number
	TotalPrice    float64             `json:"total_price"`
	PaymentStatus PaymentStatus       `json:"payment_status"` // "paid", "pending", "failed"
	OrderStatus   OrderStatus         `json:"order_status"`   // "confirmed", "processing", "shipped"
	Items         []OrderItemResponse `json:"items"`
}
