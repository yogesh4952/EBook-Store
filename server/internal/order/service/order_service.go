package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/yogesh4952/ebookstore/internal/address/models"
	bookModels "github.com/yogesh4952/ebookstore/internal/book/models"
	orderModel "github.com/yogesh4952/ebookstore/internal/order/models"
	"github.com/yogesh4952/ebookstore/internal/order/repository"
)

type IOrderServ interface {
	PlaceOrder(ctx context.Context, userID uint, orderPayload orderModel.PlaceOrderPayload) (*orderModel.PlaceOrderResponse, error)
}

type IAddressrepo interface {
	FindUserAddressById(ctx context.Context, addreessId uint) (*models.UserAddress, error)
}

type IBookRepo interface {
	FindById(ctx context.Context, id uint) (*bookModels.Book, error)
}
type orderServ struct {
	orderRepo      repository.IOrderRepo
	bookRepo       IBookRepo
	userAddresRepo IAddressrepo
}

func NewOrderService(orderRepo repository.IOrderRepo, bookRepo IBookRepo, addressRepo IAddressrepo) *orderServ {
	return &orderServ{orderRepo: orderRepo, bookRepo: bookRepo, userAddresRepo: addressRepo}
}

func (serv *orderServ) PlaceOrder(ctx context.Context, userID uint, orderPayload orderModel.PlaceOrderPayload) (*orderModel.PlaceOrderResponse, error) {

	var res orderModel.PlaceOrderResponse
	res.PaymentStatus = orderModel.PaymentStatus(orderPayload.PaymentMethod)

	userAddress, err := serv.userAddresRepo.FindUserAddressById(ctx, orderPayload.UserAddressId)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	res.Address = "City: " + userAddress.City + ", Delivery Address:" + userAddress.DeliveryAddress

	total_price := float32(0.0)
	for _, val := range orderPayload.Items {

		book, err := serv.bookRepo.FindById(ctx, val.BookId)

		if err != nil {
			return nil, fmt.Errorf("Invalid book id")
		}
		total_price += book.Price * float32(val.Quantity)
	}

	res.TotalPrice = float64(total_price)
	randOrderCode := rand.Int()

	var OrderCode string

	OrderCode = "#ORDER_CODE:" + strconv.Itoa(randOrderCode)
	res.OrderCode = OrderCode

	res.Message = "Order Placed Succesfully"
	return &res, nil
}
