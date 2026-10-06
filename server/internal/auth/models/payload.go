package models

import (
	usermodels "github.com/yogesh4952/ebookstore/internal/user/models"
)

type RegisterPayload struct {
	Firstname   string          `json:"first_name" binding:"required"`
	Lastname    string          `json:"last_name" binding:"required"`
	Email       string          `json:"email" binding:"required"`
	Role        usermodels.Role `json:"role" binding:"required" `
	Age         uint            `json:"age" binding:"required"`
	PhoneNumber string          `json:"phone_number" binding:"required"`
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type GoogleAuthRequest struct {
	Code string `json:"code" binding:"required"`
}
