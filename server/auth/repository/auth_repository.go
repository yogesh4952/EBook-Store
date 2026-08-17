package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthRepository interface {
	SaveOTP(ctx context.Context, email, otp string, duration time.Duration) error
	GetOTP(ctx context.Context, email string) (string, error)
	DeleteOTP(ctx context.Context, email string) error
}

type authRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewAuthRepository(db *gorm.DB, rdb *redis.Client) *authRepository {
	return &authRepository{db: db, rdb: rdb}
}

func (r *authRepository) SaveOTP(ctx context.Context, email, otp string, duration time.Duration) error {
	key := "otp:" + email
	return r.rdb.Set(ctx, key, otp, duration).Err()
}

func (r *authRepository) GetOTP(ctx context.Context, email string) (string, error) {
	key := "otp:" + email
	val, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("OTP expired or not requested")
	}
	return val, err
}

func (r *authRepository) DeleteOTP(ctx context.Context, email string) error {
	key := "otp:" + email
	return r.rdb.Del(ctx, key).Err()
}
