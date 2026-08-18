package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yogesh4952/ebookstore/user/models"
)

type CustomClaims struct {
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJwt(user *models.User, duration time.Duration) (string, error) {
	secretKey := os.Getenv("jwt_secret")

	bytes := []byte(secretKey)
	claims := CustomClaims{
		UserId: user.Id,
		Email:  user.Email,
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
