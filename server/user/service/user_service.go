package service

import (
	"errors"

	"github.com/yogesh4952/ebookstore/user/models"
	userrepo "github.com/yogesh4952/ebookstore/user/repository"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
}

type userService struct {
	repo userrepo.UserRepo
}

func NewUserService(repo userrepo.UserRepo) UserService {
	return &userService{repo: repo}
}

func (u *userService) GetAllUsers() ([]models.User, error) {
	users, err := u.repo.FetchAllUsers()
	if err != nil {
		return nil, errors.New("Failed to retrieve users")
	}

	if len(users) <= 0 {
		return nil, errors.New("Empty data in db")
	}

	return users, nil
}
