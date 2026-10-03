package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yogesh4952/ebookstore/internal/cart/model"
	"gorm.io/gorm"
)

type ICartrepo interface {
	AddToCart(ctx context.Context, userId uint, payload model.AddToCartPayload) (map[string]string, error)
	CartIds(ctx context.Context, userId uint) (map[string]string, error)
	GetCartQuantity(ctx context.Context, userId uint, bookId uint) (int64, error)
}

type cartRepo struct {
	db  *gorm.DB
	rdc *redis.Client
}

func NewCartRepo(db *gorm.DB, rd *redis.Client) *cartRepo {
	return &cartRepo{db: db, rdc: rd}
}

func (rp *cartRepo) AddToCart(
	ctx context.Context,
	userId uint,
	payload model.AddToCartPayload,
) (map[string]string, error) {

	cartKey := fmt.Sprintf("cart:%d", userId)
	field := fmt.Sprintf("%d", payload.BookId)

	if err := rp.rdc.HIncrBy(
		ctx,
		cartKey,
		field,
		int64(payload.Quantity),
	).Err(); err != nil {
		return nil, err
	}

	if err := rp.rdc.Expire(
		ctx,
		cartKey,
		7*24*time.Hour,
	).Err(); err != nil {
		return nil, err
	}

	cart, err := rp.rdc.HGetAll(ctx, cartKey).Result()
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (rp *cartRepo) CartIds(ctx context.Context, userId uint) (map[string]string, error) {
	key := fmt.Sprintf("cart:%d", userId)
	redishHash := rp.rdc.HGetAll(ctx, key)
	// log.Printf("%v", redishHash.Val())
	// ids := slices.Sorted(maps.Keys(redishHash.Val()))

	// book:=

	return redishHash.Val(), nil

}

func (rp *cartRepo) GetCartQuantity(ctx context.Context, userId uint, bookId uint) (int64, error) {

	key := fmt.Sprintf("cart:%d", userId)
	field := fmt.Sprintf("%d", bookId)
	quantity, err := rp.rdc.HGet(ctx, key, field).Int64()

	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	log.Printf("quantity: %v", quantity)
	return quantity, nil
}
