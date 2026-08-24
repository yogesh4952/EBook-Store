package models

import (
	bookModel "github.com/yogesh4952/ebookstore/internal/book/models"
	userModel "github.com/yogesh4952/ebookstore/internal/user/models"
	"gorm.io/gorm"
)

type Seller struct {
	gorm.Model
	UserID         uint              `json:"user_id" gorm:"not null"`
	User           userModel.User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	sellerNumber   uint              `json:"seller_number" gorm:"unique;autoIncrement"`
	PublishedBooks []*bookModel.Book `gorm:"foreignKey:SellerId"`
}
