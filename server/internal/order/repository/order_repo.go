package repository

import "gorm.io/gorm"

type IOrderRepo interface {
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *orderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) PlaceOrder() error {
	return nil
}
