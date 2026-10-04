package models

import (
	"github.com/yogesh4952/ebookstore/internal/user/models"
	baseModels "github.com/yogesh4952/ebookstore/pkg/models"
)

type UserAddress struct {
	baseModels.Base
	City            string      `json:"city" binding:"required" gorm:"not null"`
	DeliveryAddress string      `json:"delivery_address" binding:"required" gorm:"not null"`
	UserId          uint        `json:"user_id" binding:"omitempty"`
	User            models.User `json:"user" gorm:"foreignKey:UserId"`
}
