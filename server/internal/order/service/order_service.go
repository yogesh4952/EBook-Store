package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"

	"github.com/yogesh4952/ebookstore/internal/address/models"
	bookModels "github.com/yogesh4952/ebookstore/internal/book/models"
	orderModel "github.com/yogesh4952/ebookstore/internal/order/models"
	"github.com/yogesh4952/ebookstore/internal/order/repository"
	"gorm.io/gorm"
)

var (
	ErrInvalidAddress       = errors.New("invalid address id")
	ErrBookNotFound         = errors.New("invalid book id")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrPlacingOrder         = errors.New("failed to place order")
)

type IOrderServ interface {
	PlaceOrder(ctx context.Context, userID uint, orderPayload orderModel.PlaceOrderPayload) (*orderModel.PlaceOrderResponse, error)
}

type IAddressrepo interface {
	FindUserAddressByIdAndUser(ctx context.Context, addressID uint, userID uint) (*models.UserAddress, error)
}

type IBookRepo interface {
	FindByIds(ctx context.Context, ids []uint) ([]bookModels.Book, error)
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
	switch orderPayload.PaymentMethod {
	case orderModel.PayementCOD, orderModel.PayementEsewa:
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidPaymentMethod, orderPayload.PaymentMethod)
	}

	userAddress, err := serv.userAddresRepo.FindUserAddressByIdAndUser(ctx, orderPayload.UserAddressId, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %d", ErrInvalidAddress, orderPayload.UserAddressId)
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidAddress, err)
	}

	bookIDs := make([]uint, 0, len(orderPayload.Items))

	for _, item := range orderPayload.Items {
		bookIDs = append(bookIDs, item.BookId)
	}

	books, err := serv.bookRepo.FindByIds(ctx, bookIDs)
	if err != nil {
		return nil, fmt.Errorf("Failed to load books: %w", err)
	}

	bookMap := make(map[uint]*bookModels.Book, len(books))
	for i := range books {
		bookMap[books[i].ID] = &books[i]
	}

	quantities := make(map[uint]uint, len(orderPayload.Items))
	for _, item := range orderPayload.Items {
		if _, ok := bookMap[item.BookId]; !ok {
			return nil, fmt.Errorf("%w: %d", ErrBookNotFound, item.BookId)
		}
		quantities[item.BookId] += item.Quantity
	}

	var totalPrice int64
	items := make([]*orderModel.OrderItem, 0, len(quantities))
	for bookID, quantity := range quantities {
		book := bookMap[bookID]
		totalPrice += book.Price * int64(quantity)
		items = append(items, &orderModel.OrderItem{
			BookId:    bookID,
			Quantity:  quantity,
			UnitPrice: book.Price,
		})
	}

	orderCode, err := generateOrderCode()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPlacingOrder, err)
	}

	data := orderModel.Order{
		OrderCode:               orderCode,
		OrderStatus:             orderModel.OrderPlaced,
		UserId:                  userID,
		PaymentMethod:           orderPayload.PaymentMethod,
		ShippingCity:            userAddress.City,
		ShippingDeliveryAddress: userAddress.DeliveryAddress,
		TotalPrice:              totalPrice,
	}

	if err := serv.orderRepo.PlaceOrder(ctx, &data, items); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPlacingOrder, err)
	}

	res := orderModel.PlaceOrderResponse{
		OrderID:       data.ID,
		OrderCode:     data.OrderCode,
		TotalPrice:    data.TotalPrice,
		PaymentStatus: data.PaymentStatus,
		OrderStatus:   data.OrderStatus,
		Address:       data.ShippingCity + "," + data.ShippingDeliveryAddress,
		Items:         make([]orderModel.OrderItemResponse, 0, len(items)),
	}

	for _, item := range items {
		book := bookMap[item.BookId]
		res.Items = append(res.Items, orderModel.OrderItemResponse{
			BookId:    item.BookId,
			Title:     book.Title,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.UnitPrice * int64(item.Quantity),
		})
	}

	return &res, nil
}

func generateOrderCode() (string, error) {
	raw := make([]byte, 6)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return "ORD-" + base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw), nil
}
