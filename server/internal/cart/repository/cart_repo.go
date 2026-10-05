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
	ClearCart(c context.Context, userId uint) error
	RemoveFromCart(ctx context.Context, userId uint, payload model.AddToCartPayload) (map[string]string, error)
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
		return nil, fmt.Errorf("failed to increment cart item: %w", err)
	}

	if err := rp.rdc.Expire(ctx, cartKey, 7*24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("failed to refresh cart expiration: %w", err)
	}

	return rp.rdc.HGetAll(ctx, cartKey).Result()
}

func (rp *cartRepo) RemoveFromCart(ctx context.Context, userId uint, payload model.AddToCartPayload) (map[string]string, error) {
	cartKey := fmt.Sprintf("cart:%d", userId)
	field := fmt.Sprintf("%d", payload.BookId)

	current, err := rp.rdc.HGet(ctx, cartKey, field).Int64()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("book %d is not in the cart", payload.BookId)
	}
	if err != nil {
		return nil, err
	}

	remaining := current - int64(payload.Quantity)
	if remaining < 0 {
		return nil, fmt.Errorf("cannot remove %d, only %d in cart", payload.Quantity, current)
	}

	if remaining == 0 {
		if err := rp.rdc.HDel(ctx, cartKey, field).Err(); err != nil {
			return nil, err
		}
	} else if err := rp.rdc.HIncrBy(ctx, cartKey, field, -int64(payload.Quantity)).Err(); err != nil {
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

func (rp *cartRepo) ClearCart(c context.Context, userId uint) error {
	key := fmt.Sprintf("cart:%d", userId)
	if err := rp.rdc.Del(c, key).Err(); err != nil {
		return err
	}
	return nil
}
