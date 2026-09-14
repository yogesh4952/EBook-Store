package models

import (
	"errors"

	"github.com/yogesh4952/ebookstore/internal/user/models"
	"gorm.io/gorm"
)

type (
	PaymentMethod string
	PaymentStatus string
	OrderStatus   string
)

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
	OrderPlaced       OrderStatus = "PLACED"
)

type Order struct {
	gorm.Model
	OrderCode       string        `json:"order_code" gorm:"unique;not null"`
	TransactionUUID string        `gorm:"unique;" json:"transaction_uuid"`
	PaymentMethod   PaymentMethod `json:"payment_method" gorm:"check:payment_method IN ('COD','ESEWA');not null;"`
	PaymentStatus   PaymentStatus `json:"payment_status" gorm:"check:payment_status IN ('PAID','PENDING','REFUND');not null; default:'PENDING'"`
	OrderStatus     OrderStatus   `json:"order_status" gorm:"check:order_status IN ('DELIVERED','CANCELLED' ,'PLACED'); default:''"`
	TotalPrice      float32       `json:"total_price" binding:"required,min=1"`
	UserId          uint          `json:"user_id" gorm:"index:idx_userid"`

	ShippingCity            string `json:"shipping_city" binding:"required" gorm:"not null"`
	ShippingDeliveryAddress string `json:"shipping_delivery_address" binding:"required" gorm:"not null"`

	User      models.User `gorm:"foreignKey:UserId"`
	OrderItem []OrderItem `json:"order_items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
	OrderID       uint                `json:"order_id"`
	OrderCode     string              `json:"order_code"` // Human-readable order number
	TotalPrice    float32             `json:"total_price"`
	PaymentStatus PaymentStatus       `json:"payment_status"` // "paid", "pending", "failed"
	OrderStatus   OrderStatus         `json:"order_status"`   // "confirmed", "processing", "shipped"
	Address       string              `json:"address"`
	Items         []OrderItemResponse `json:"items"`
	EsewaPayload  *EsewaPayload       `json:"esewa_payload,omitempty"`
}

// "amount": "100",
// "failure_url": "https://developer.esewa.com.np/failure",
// "product_delivery_charge": "0",
// "product_service_charge": "0",
// "product_code": "EPAYTEST",
// "signature": "i94zsd3oXF6ZsSr/kGqT4sSzYQzjj1W/waxjWyRwaME=",
// "signed_field_names": "total_amount,transaction_uuid,product_code",
// "success_url": "https://developer.esewa.com.np/success",
// "tax_amount": "10",
// "total_amount": "110",
// "transaction_uuid": "241028"

type EsewaPayload struct {
	Amount                string `json:"amount"`
	TotalAmount           string `json:"total_amount"`
	TransactionUUID       string `json:"transaction_uuid"`
	ProductCode           string `json:"product_code"`
	ProductServiceCharge  string `json:"product_service_charge"`
	ProductDeliveryCharge string `json:"product_delivery_charge"`
	TaxAmount             string `json:"tax_amount"`
	SignedFieldNames      string `json:"signed_field_names"`
	Signature             string `json:"signature"`
	SuccessUrl            string `json:"success_url"`
	FailureUrl            string `json:"failure_url"`
}

var (
	ErrInvalidAddress       = errors.New("invalid address id")
	ErrBookNotFound         = errors.New("invalid book id")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrPlacingOrder         = errors.New("failed to place order")
)
