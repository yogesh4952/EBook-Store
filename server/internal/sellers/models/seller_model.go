package models

import (
	bookModel "github.com/yogesh4952/ebookstore/internal/book/models"
	"gorm.io/gorm"
)

type Seller struct {
	gorm.Model
	VendorNumber   uint             `json:"vendor_number" binding:"required"`
	PublishedBooks []bookModel.Book `gorm:"foreignKey:ID"`
}
