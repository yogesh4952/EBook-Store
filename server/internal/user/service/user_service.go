package service

import (
	"context"
	"errors"

	"github.com/yogesh4952/ebookstore/internal/user/models"
	userrepo "github.com/yogesh4952/ebookstore/internal/user/repository"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type userService struct {
	repo userrepo.IUser
}

func NewUserService(repo userrepo.IUser) UserService {
	return &userService{repo: repo}
}

func (u *userService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := u.repo.FetchAllUsers(ctx)
	if err != nil {
		return nil, errors.New("Failed to retrieve users")
	}

	if len(users) <= 0 {
		return nil, errors.New("Empty data in db")
	}

	return users, nil
}
