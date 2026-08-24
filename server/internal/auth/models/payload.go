package models

import (
	usermodels "github.com/yogesh4952/ebookstore/internal/user/models"
)

type RegisterPayload struct {
	Firstname string          `json:"first_name" binding:"required"`
	Lastname  string          `json:"last_name" binding:"required"`
	Email     string          `json:"email" binding:"required"`
	Role      usermodels.Role `json:"role" binding:"required" `
	Age       uint            `json:"age" binding:"required"`
}
