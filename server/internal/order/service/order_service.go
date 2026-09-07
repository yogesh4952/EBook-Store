package service

import (
	"context"

	orderModel "github.com/yogesh4952/ebookstore/internal/order/models"
	"github.com/yogesh4952/ebookstore/internal/order/repository"
)

type IOrderServ interface{}

type orderServ struct {
	orderRepo repository.IOrderRepo
}

func NewOrderService(orderRepo repository.IOrderRepo) *orderServ {
	return &orderServ{orderRepo: orderRepo}
}

func (serv *orderServ) PlaceOrder(ctx context.Context, data *orderModel.Order) error {
	return nil
}
