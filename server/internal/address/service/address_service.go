package service

import (
	"context"

	"github.com/yogesh4952/ebookstore/internal/address/models"
	"github.com/yogesh4952/ebookstore/internal/address/repository"
	userRepo "github.com/yogesh4952/ebookstore/internal/user/repository"
)

type IAddressService interface{}

type addressService struct {
	addressRepo repository.IAddressrepo
	userRepo    userRepo.IUser
}

func NewAddressService(addressRepo repository.IAddressrepo,
	userRepo userRepo.IUser) *addressService {
	return &addressService{
		addressRepo: addressRepo,
		userRepo:    userRepo,
	}
}

func (as *addressService) AddAddress(ctx context.Context, payload *models.UserAddress) error {
	// userId := payload.UserId
	return nil

	// data, err := as.userRepo.()

}
