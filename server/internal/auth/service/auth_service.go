package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	Google(ctx context.Context, code string) (string, error)
}

type UserLookup interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type TokenGenerator interface {
	GenerateJwt(user *models.User, duration time.Duration) (string, error)
}

type EmailSender interface {
	SentOTPEmail(toEmail, otp string) error
	GenerateOTP() (string, error)
}

type authService struct {
	repo         repository.IAuthRepository
	userStore    UserLookup
	tokenGen     TokenGenerator
	emailService EmailSender
	client       *http.Client
}

func NewAuthService(
	repo repository.IAuthRepository, userStore UserLookup,
	tokenGen TokenGenerator,
	emailService EmailSender,
) *authService {
	return &authService{
		repo:         repo,
		userStore:    userStore,
		tokenGen:     tokenGen,
		emailService: emailService,
		client:       &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *authService) SendOTP(ctx context.Context, email string) error {

	if _, err := s.userStore.FindByEmail(ctx, email); err != nil {
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

func (s *authService) VerifyOTP(ctx context.Context, email, inputOTP string) (bool, error) {
	storedOTP, err := s.repo.GetOTP(ctx, email)
	if err != nil {
		return false, err
	}

	if storedOTP != inputOTP {
		return false, errors.New("invalid verification code")
	}

	_ = s.repo.DeleteOTP(ctx, email)

	return true, nil
}

func (s *authService) Login(ctx context.Context, email, inputOTP string) (string, error) {
	ok, err := s.VerifyOTP(ctx, email, inputOTP)
	if !ok || err != nil {
		return "", errors.New("invalid or expired verification code")
	}

	user, err := s.userStore.FindByEmail(ctx, email)
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

func (s *authService) Google(ctx context.Context, code string) (string, error) {
	base := os.Getenv("GOOGLE_URL")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	log.Printf("[GOOGLE] base=%q redirectURI=%q clientID=%q secret_len=%d", base, redirectURI, clientID, len(clientSecret))
	if base == "" || redirectURI == "" || clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("google credentials not set")
	}

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("google server unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return "", fmt.Errorf("google token exchange failed with status %d: %v", resp.StatusCode, errBody)
	}

	var tokenResp authmodels.GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode google response: %w", err)
	}
	log.Printf("TOKEN: %s", tokenResp.IDToken)

	parts := strings.Split(tokenResp.IDToken, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid google id_token")
	}
	payloadB64 := parts[1]
	switch len(payloadB64) % 4 {
	case 2:
		payloadB64 += "=="
	case 3:
		payloadB64 += "="
	}
	payloadBytes, err := base64.URLEncoding.DecodeString(payloadB64)
	if err != nil {
		return "", fmt.Errorf("failed to decode google id_token payload: %w", err)
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return "", fmt.Errorf("failed to unmarshal google claims: %w", err)
	}
	email, _ := claims["email"].(string)
	firstname, _ := claims["first_name"].(string)
	lastname, _ := claims["last_name"].(string)
	if email == "" {
		return "", fmt.Errorf("google id_token missing email")
	}
	if firstname == "" {
		firstname = "Google"
	}
	if lastname == "" {
		lastname = "User"
	}
	user, err := s.userStore.FindByEmail(ctx, email)
	if err != nil {
		userRegistration := &models.User{
			Firstname:   firstname,
			Lastname:    lastname,
			Email:       email,
			Role:        models.RoleCustomer,
			PhoneNumber: "0000000000",
		}
		if regErr := s.repo.RegisterUser(ctx, userRegistration); regErr != nil {
			return "", fmt.Errorf("Failed to register user: %w", regErr)
		}
		user, err = s.userStore.FindByEmail(ctx, email)
		if err != nil {
			return "", fmt.Errorf("Failed to find registered user: %w", err)
		}
	}

	jwtToken, err := s.tokenGen.GenerateJwt(user, time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate internal jwt: %w", err)
	}
	return jwtToken, nil
}
