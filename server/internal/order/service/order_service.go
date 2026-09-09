package service

import (
	"context"
	"encoding/json"
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

	var data orderModel.Order
	data.OrderStatus = orderModel.OrderStatus("PLACED")
	data.UserId = userID

	userAddress, err := serv.userAddresRepo.FindUserAddressById(ctx, orderPayload.UserAddressId)
	data.PaymentMethod = orderPayload.PaymentMethod

	data.User = userAddress.User
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	data.ShippingCity = userAddress.City
	data.ShippingDeliveryAddress = userAddress.DeliveryAddress

	total_price := float32(0.0)
	for _, val := range orderPayload.Items {

		book, err := serv.bookRepo.FindById(ctx, val.BookId)

		// we have to push the data into order item table as well
		if err != nil {
			return nil, fmt.Errorf("Invalid book id")
		}
		total_price += book.Price * float32(val.Quantity)
	}

	data.TotalPrice = float32(total_price)
	randOrderCode := rand.Int()

	var OrderCode string

	OrderCode = "#ORDER_CODE:" + strconv.Itoa(randOrderCode)
	data.OrderCode = OrderCode

	err = serv.orderRepo.PlaceOrder(ctx, &data)

	if err != nil {
		return nil, err
	}

	var res orderModel.PlaceOrderResponse
	res.Address = data.ShippingCity + "," + data.ShippingDeliveryAddress
	res.TotalPrice = float64(data.TotalPrice)
	res.OrderCode = data.OrderCode
	items, err := json.Marshal(orderPayload.Items)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(items, &res.Items); err != nil {
		return nil, err
	}
	res.OrderStatus = "PLACED"
	res.PaymentStatus = data.PaymentStatus
	res.OrderID = data.ID

	return &res, err
}
