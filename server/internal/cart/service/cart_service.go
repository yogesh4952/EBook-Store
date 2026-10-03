package service

import (
	"context"
	"fmt"

	bookModels "github.com/yogesh4952/ebookstore/internal/book/models"
	"github.com/yogesh4952/ebookstore/internal/cart/model"
)

type ICartService interface {
	AddToCart(ctx context.Context, userId uint, payload model.AddToCartPayload) (map[string]string, error)
	GetCartItem(ctx context.Context, userId uint) (*model.CartResponse, error)
}

type ICartRepo interface {
	GetCartQuantity(ctx context.Context, userId uint, bookId uint) (int64, error)
	AddToCart(ctx context.Context, userId uint, payload model.AddToCartPayload) (map[string]string, error)
	CartList(ctx context.Context, userId uint) (*model.CartResponse, error)
}

type cartService struct {
	cartRepo ICartRepo
	bookRepo IBookLookup
}

// IBookLookup is the narrow slice of the book repository this service needs.
// Declaring it here keeps the cart service from importing the book repository
// package directly.
type IBookLookup interface {
	FindById(ctx context.Context, id uint) (*bookModels.Book, error)
}

func NewCartService(cartRepo ICartRepo, br IBookLookup) *cartService {
	return &cartService{cartRepo: cartRepo, bookRepo: br}
}

func (cs *cartService) AddToCart(ctx context.Context, userId uint, payload model.AddToCartPayload) (map[string]string, error) {
	book, err := cs.bookRepo.FindById(ctx, payload.BookId)
	if err != nil {
		return nil, fmt.Errorf("book %d not found: %w", payload.BookId, err)
	}

	redisQuantity, err := cs.cartRepo.GetCartQuantity(ctx, userId, payload.BookId)
	if err != nil {
		return nil, err
	}

	if (book.Units < payload.Quantity) || (redisQuantity+int64(payload.Quantity)) > int64(book.Units) {
		return nil, fmt.Errorf("only %d units of book %d available", book.Units, payload.BookId)
	}

	return cs.cartRepo.AddToCart(ctx, userId, payload)
}

func (cs *cartService) GetCartItem(ctx context.Context, userId uint) (*model.CartResponse, error) {
	return cs.cartRepo.CartList(ctx, userId)
}
