package userrepo

import (
	"errors"

	"github.com/yogesh4952/ebookstore/user/models"
	"gorm.io/gorm"
)

type UserRepo interface {
	FindByEmail(email string) (*models.User, error)
	FetchAllUsers() ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := u.db.Where("email = ? ", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}
func (u *userRepository) FetchAllUsers() ([]models.User, error) {
	var users []models.User

	err := u.db.Find(&users).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.User{}, nil
		}
		return nil, err
	}
	return users, nil

}
