package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yogesh4952/ebookstore/internal/auth"
	usermodels "github.com/yogesh4952/ebookstore/internal/user/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	SaveOTP(ctx context.Context, email, otp string, duration time.Duration) error
	GetOTP(ctx context.Context, email string) (string, error)
	DeleteOTP(ctx context.Context, email string) error
	RegisterUser(ctx context.Context, user *usermodels.User) error
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

// auth_repository.go
func (r *authRepository) RegisterUser(ctx context.Context, user *usermodels.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return auth.ErrDuplicateEmail

	}

	return fmt.Errorf("failed to create user in db: %w", err)
}
