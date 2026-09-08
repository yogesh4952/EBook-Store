package service

import (
	"context"
	"fmt"
	"math/rand"

	adddressRepo "github.com/yogesh4952/ebookstore/internal/address/repository"
	bookRepo "github.com/yogesh4952/ebookstore/internal/book/repository"
	orderModel "github.com/yogesh4952/ebookstore/internal/order/models"
	"github.com/yogesh4952/ebookstore/internal/order/repository"
)

type IOrderServ interface {
	PlaceOrder(ctx context.Context, userID uint, orderPayload orderModel.PlaceOrderPayload) (*orderModel.PlaceOrderResponse, error)
}

type orderServ struct {
	orderRepo      repository.IOrderRepo
	bookRepo       bookRepo.IBookRepo
	userAddresRepo adddressRepo.IAddressrepo
}

func NewOrderService(orderRepo repository.IOrderRepo, bookRepo bookRepo.IBookRepo, addressRepo adddressRepo.IAddressrepo) *orderServ {
	return &orderServ{orderRepo: orderRepo, bookRepo: bookRepo, userAddresRepo: addressRepo}
}

func (serv *orderServ) PlaceOrder(ctx context.Context, userID uint, orderPayload orderModel.PlaceOrderPayload) (*orderModel.PlaceOrderResponse, error) {

	var res orderModel.PlaceOrderResponse
	res.PaymentStatus = orderModel.PaymentStatus(orderPayload.PaymentMethod)

	userAddress, err := serv.userAddresRepo.FindUserAddressById(ctx, orderPayload.UserAddressId)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("Invalid user address")
	}

	res.Address = "City: " + userAddress.City + ", Delivery Address:" + userAddress.DeliveryAddress

	total_price := float32(0.0)
	for _, val := range orderPayload.Items {

		book, err := serv.bookRepo.FindById(ctx, val.BookId)

		if err != nil {
			return nil, fmt.Errorf("Invalid book id")
		}
		total_price = (float32(total_price) + book.Price) * float32(val.Quantity)
	}

	res.TotalPrice = float64(total_price)
	randOrderCode := rand.Int()

	var OrderCode string

	OrderCode = "#ORDER_CODE:" + string(randOrderCode)
	res.OrderCode = OrderCode

	res.Message = "Order Placed Succesfully"
	return &res, nil
}
