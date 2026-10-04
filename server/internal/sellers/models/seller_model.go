package models

import (
	userModel "github.com/yogesh4952/ebookstore/internal/user/models"
	baseModels "github.com/yogesh4952/ebookstore/pkg/models"
)

type Seller struct {
	baseModels.Base

	UserID       uint           `json:"user_id" gorm:"not null"`
	User         userModel.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SellerNumber uint           `json:"seller_number" gorm:"unique"`
}
