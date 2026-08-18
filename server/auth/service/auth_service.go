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
	VerifyOTP(ctx context.Context, email, inputOTP string) (*models.User, string, error)
	Login(ctx context.Context, email string) (string, error)
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

func (s *authService) VerifyOTP(ctx context.Context, email, inputOTP string) (*models.User, string, error) {
	storedOTP, err := s.repo.GetOTP(ctx, email)
	if err != nil {
		return nil, "", err
	}

	// 2. Compare inputs
	if storedOTP != inputOTP {
		return nil, "", errors.New("invalid verification code")
	}

	// 3. Delete from Redis immediately (Single-use enforcement)
	_ = s.repo.DeleteOTP(ctx, email)

	// 5. Generate authentication token placeholder
	token := "sample_jwt_access_token"

	// Create a user object. At this stage we only have the email from the
	// verification step; further user data can be fetched or created as
	// needed by the caller or additional service logic.
	user := &models.User{Email: email}

	return user, token, nil
}

func (s *authService) Login(ctx context.Context, email string) (string, error) {
	// Use the exported lookup function from the user repository package.
	// The concrete repository exposes an exported helper FindByEmail.
	user, err := s.userStore.FindByEmail(email)

	if err != nil {
		return "", fmt.Errorf("error fetching user with email %s: %w", email, err)
	}

	token, err := utils.GenerateJwt(user, time.Hour)

	if err != nil {
		return "", fmt.Errorf("Error creating jwt token: %w", err)
	}

	// TODO: generate and return a real JWT/token. Returning placeholder for now.
	return token, nil
}
