package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yogesh4952/ebookstore/internal/user/models"
)

type JwtManager struct{}

func NewJwtManager() *JwtManager {
	return &JwtManager{}
}

type CustomClaims struct {
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JwtManager) GenerateJwt(user *models.User, duration time.Duration) (string, error) {
	secretKey := os.Getenv("jwt_secret")

	bytes := []byte(secretKey)
	claims := CustomClaims{
		UserId: user.ID,
		Email:  user.Email,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "my-backend-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(bytes)

	if err != nil {
		return "", err

	}

	return signedToken, nil
}

func ValidateJwt(tokenString string) (*CustomClaims, error) {
	secretKey := os.Getenv("jwt_secret")
	bytes := []byte(secretKey)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Method.Alg())
		}
		return bytes, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer("my-backend-service"),
	)

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
