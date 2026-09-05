package models

import (
	"github.com/yogesh4952/ebookstore/internal/user/models"
	"gorm.io/gorm"
)

type UserAddresses struct {
	gorm.Model
	City            string      `json:"city" binding:"required" gorm:"not null"`
	DeliveryAddress string      `json:"delivery_address" binding:"required" gorm:"not null"`
	PhoneNumber     string      `json:"phone_number" binding:"required,min=10,max=10"`
	UserId          uint        `json:"user_id"`
	User            models.User `json:"user" gorm:"foreignKey:UserId"`
}
