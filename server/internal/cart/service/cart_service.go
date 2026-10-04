package service

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/yogesh4952/ebookstore/internal/book/models"
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
	CartIds(ctx context.Context, userId uint) (map[string]string, error)
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
	FindByIds(ctx context.Context, ids []uint) ([]models.Book, error)
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
	raw, err := cs.cartRepo.CartIds(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("Error during fetching cart items :%v", err)
	}

	out := make(map[uint]uint, len(raw))

	for k, v := range raw {
		bookId, err := strconv.ParseUint(k, 10, 64)
		if err != nil {
			return nil, err
		}

		qty, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}

		out[uint(bookId)] = uint(qty)

	}

	ids := make([]uint, 0, len(out))

	for id := range out {
		ids = append(ids, id)

	}
	slices.Sort(ids)
	bookData, err := cs.bookRepo.FindByIds(ctx, ids)

	if err != nil {
		return nil, fmt.Errorf("Error fetching books: %v", err)
	}

	items := make([]model.ItemResp, 0, len(bookData))

	var subtotal float32
	subtotal = 0

	var total float32
	for _, b := range bookData {
		subtotal = b.Price * float32(out[b.ID])

		items = append(items, model.ItemResp{
			Data:     b,
			Quantity: out[b.ID],
			Subtotal: subtotal,
		})

		total += subtotal

	}

	return &model.CartResponse{Items: items, Total: total}, nil
}
