package userrepo

import (
	"github.com/yogesh4952/ebookstore/models"
	"gorm.io/gorm"
)

type UserRepo interface {
	findByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (u *userRepository) findByEmail(email string) (*models.User, error) {
	var user models.User

	err := u.db.Where("email = ? ", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}
