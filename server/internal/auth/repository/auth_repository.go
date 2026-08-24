package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yogesh4952/ebookstore/internal/auth"
	sellerModel "github.com/yogesh4952/ebookstore/internal/sellers/models"
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

func (r *authRepository) RegisterUser(ctx context.Context, user *usermodels.User) error {
	// Execute both operations inside an isolated DB transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return auth.ErrDuplicateEmail
			}
			return err
		}

		if user.Role == usermodels.Roleseller {
			seller := &sellerModel.Seller{
				UserID: user.ID,
			}

			if err := tx.Create(seller).Error; err != nil {
				return fmt.Errorf("failed to create seller profile: %w", err)
			}
		}

		return nil
	})
}
