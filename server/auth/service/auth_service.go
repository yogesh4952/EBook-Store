package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yogesh4952/ebookstore/auth/repository"
	"github.com/yogesh4952/ebookstore/pkg/utils"
	"github.com/yogesh4952/ebookstore/user/models"
)

type AuthService interface {
	SendOTP(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, email, inputOTP string) (bool, error)
	Login(ctx context.Context, email, inputotp string) (string, error)
}

type UserLookup interface {
	FindByEmail(email string) (*models.User, error)
}

type authService struct {
	repo      repository.AuthRepository
	userStore UserLookup
}

func NewAuthService(repo repository.AuthRepository, userStore UserLookup) AuthService {
	return &authService{repo: repo, userStore: userStore}
}

// SendOTP is a func that sent otp to the user and store in redis
func (s *authService) SendOTP(ctx context.Context, email string) error {
	otp, err := utils.GenerateOTP()

	if err != nil {
		return errors.New("Failed to generate verification code")
	}

	if err := s.repo.SaveOTP(ctx, email, otp, 5*time.Minute); err != nil {
		return errors.New("Failed to save verification session")
	}

	go func() {
		_ = utils.SentOTPEmail(email, otp)
	}()
	return nil

}

// VerifyOTP is a func that is used to compare the otp from user written and the otp store in the redis

func (s *authService) VerifyOTP(ctx context.Context, email, inputOTP string) (bool, error) {
	storedOTP, err := s.repo.GetOTP(ctx, email)
	if err != nil {
		return false, err
	}

	// 2. Compare inputs
	if storedOTP != inputOTP {
		return false, errors.New("invalid verification code")
	}

	// 3. Delete from Redis immediately (Single-use enforcement)
	_ = s.repo.DeleteOTP(ctx, email)

	return true, nil
}

func (s *authService) Login(ctx context.Context, email, inputOTP string) (string, error) {
	storedOTP, err := s.repo.GetOTP(ctx, email)
	if err != nil {
		return "", errors.New("verification code expired or not requested")
	}

	if storedOTP != inputOTP {
		return "", errors.New("invalid verification code")
	}

	_ = s.repo.DeleteOTP(ctx, email)

	user, err := s.userStore.FindByEmail(email)
	if err != nil {
		return "", fmt.Errorf("user not found for email %s: %w", email, err)
	}

	token, err := utils.GenerateJwt(user, time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}
