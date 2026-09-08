package userrepo

import (
	"context"
	"errors"

	"github.com/yogesh4952/ebookstore/internal/user/models"
	"gorm.io/gorm"
)

type IUser interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FetchAllUsers(ctx context.Context) ([]models.User, error)
	FindById(ctx context.Context, id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := u.db.WithContext(ctx).Where("email = ? ", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (u *userRepository) FindById(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	err := u.db.First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) FetchAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := u.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.User{}, nil
		}
		return nil, err
	}
	return users, nil

}
