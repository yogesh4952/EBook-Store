package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yogesh4952/ebookstore/internal/auth"
	authmodels "github.com/yogesh4952/ebookstore/internal/auth/models"
	"github.com/yogesh4952/ebookstore/internal/auth/repository"
	"github.com/yogesh4952/ebookstore/internal/user/models"
)

type AuthService interface {
	SendOTP(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, email, inputOTP string) (bool, error)
	Login(ctx context.Context, email, inputotp string) (string, error)
	Register(ctx context.Context, data *authmodels.RegisterPayload) (string, error)
}

type UserLookup interface {
	FindByEmail(email string) (*models.User, error)
}

type TokenGenerator interface {
	GenerateJwt(user *models.User, duration time.Duration) (string, error)
}

type EmailSender interface {
	SentOTPEmail(toEmail, otp string) error
	GenerateOTP() (string, error)
}

type authService struct {
	repo         repository.AuthRepository
	userStore    UserLookup
	tokenGen     TokenGenerator
	emailService EmailSender
}

func NewAuthService(
	repo repository.AuthRepository,
	userStore UserLookup,
	tokenGen TokenGenerator,
	emailService EmailSender,
) *authService {
	return &authService{
		repo:         repo,
		userStore:    userStore,
		tokenGen:     tokenGen,
		emailService: emailService,
	}
}

// SendOTP is a func that sent otp to the user and store in redis
func (s *authService) SendOTP(ctx context.Context, email string) error {

	if _, err := s.userStore.FindByEmail(email); err != nil {
		return errors.New("invalid email")
	}

	otp, err := s.emailService.GenerateOTP()

	if err != nil {
		return errors.New("Failed to generate verification code")
	}

	if err := s.repo.SaveOTP(ctx, email, otp, 5*time.Minute); err != nil {
		return errors.New("Failed to save verification session")
	}

	return s.emailService.SentOTPEmail(email, otp)

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
	ok, err := s.VerifyOTP(ctx, email, inputOTP)
	if !ok || err != nil {
		return "", errors.New("invalid or expired verification code")
	}

	user, err := s.userStore.FindByEmail(email)
	if err != nil {
		return "", fmt.Errorf("user not found for email %s: %w", email, err)
	}

	return s.tokenGen.GenerateJwt(user, time.Hour)
}

func (s *authService) Register(ctx context.Context, payload *authmodels.RegisterPayload) (string, error) {

	if !payload.Role.IsValid() {

		return "", auth.ErrInvalid
	}
	user := &models.User{
		Firstname: payload.Firstname,
		Lastname:  payload.Lastname,
		Email:     payload.Email,
		Role:      payload.Role,
		Age:       payload.Age,
	}
	err := s.repo.RegisterUser(ctx, user)

	if err != nil {
		if errors.Is(err, auth.ErrDuplicateEmail) {

			return "", auth.ErrDuplicateEmail
		}

		return "", fmt.Errorf("registration failed: %w", err)
	}

	return "User registered successfully", nil
}
