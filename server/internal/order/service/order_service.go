package service

import (
	"context"
	"fmt"
	"math/rand"

	orderModel "github.com/yogesh4952/ebookstore/internal/order/models"
	"github.com/yogesh4952/ebookstore/internal/order/repository"
)

type IOrderServ interface {
	PlaceOrder(ctx context.Context, userID uint, orderPayload orderModel.PlaceOrderPayload) (*orderModel.PlaceOrderResponse, error)
}

type orderServ struct {
	orderRepo repository.IOrderRepo
}

func NewOrderService(orderRepo repository.IOrderRepo) *orderServ {
	return &orderServ{orderRepo: orderRepo}
}

func (serv *orderServ) PlaceOrder(ctx context.Context, userID uint, orderPayload orderModel.PlaceOrderPayload) (*orderModel.PlaceOrderResponse, error) {

	var res orderModel.PlaceOrderResponse
	// order code
	res.PaymentStatus = orderModel.PaymentStatus(orderPayload.PaymentMethod)

	userAddress, err = serv.userAddressRepo.findUserAddressById(orderPayload.UserAddressId)

	if err != nil {
		return nil, fmt.Errorf("Invalid user address")
	}

	res.Address = userAddress
	//  item array bhitra rahekoo book id validate garna paro

	total_price := 0
	for _, val := range orderPayload.Items {

		book, err := serv.bookRepo.findBookById(val.BookId)

		if err != nil {
			return nil, fmt.Errorf("Invalid book id")
		}
		total_price = (price + book.Price) * int(val.Quantity)
	}

	res.TotalPrice = float64(total_price)
	randOrderCode := rand.Int()

	OrderCode := "#ORDER_CODE" + randOrderCode
	res.OrderCode = OrderCode

	res.Message = "Order Placed Succesfully"
	return &res, nil
}
