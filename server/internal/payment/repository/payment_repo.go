package repository

import "gorm.io/gorm"

type IPaymentRepo interface {
}

type paymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *paymentRepo {
	return &paymentRepo{db: db}
}
