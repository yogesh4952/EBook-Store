package service

import (
	"context"

	"github.com/yogesh4952/ebookstore/internal/address/models"
	"github.com/yogesh4952/ebookstore/internal/address/repository"
	usermodels "github.com/yogesh4952/ebookstore/internal/user/models"
)

type UserMethod interface {
	FindById(ctx context.Context, id uint) (*usermodels.User, error)
}
type IAddressService interface {
	AddAddress(ctx context.Context, payload *models.UserAddress) error
}

type addressService struct {
	addressRepo repository.IAddressrepo
	userRepo    UserMethod
}

func NewAddressService(addressRepo repository.IAddressrepo,
	userRepo UserMethod) *addressService {
	return &addressService{
		addressRepo: addressRepo,
		userRepo:    userRepo,
	}
}

func (as *addressService) AddAddress(ctx context.Context, payload *models.UserAddress) error {
	userId := payload.UserId

	_, err := as.userRepo.FindById(ctx, userId)

	if err != nil {
		return err
	}

	err = as.addressRepo.AddAddress(ctx, payload)
	if err != nil {
		return err
	}

	return nil

	// data, err := as.userRepo.()

}
