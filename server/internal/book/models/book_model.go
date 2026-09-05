package models

import (
	"github.com/yogesh4952/ebookstore/internal/sellers/models"
	"gorm.io/gorm"
)

type Status string

const (
	StatusAvailable   Status = "in_stock"
	StatusUnavailable Status = "out_of_stock"
)

type Book struct {
	gorm.Model
	Title        string         `json:"title"`
	AuthorName   string         `json:"author_name"`
	Genre        string         `json:"genre"`
	Category     string         `json:"category"`
	Pages        uint           `json:"pages"`
	Publication  string         `json:"publication"`
	Price        float32        `json:"price"`
	Units        int            `json:"units"`
	CoverPageUrl string         `json:"cover_page_url"`
	SellerID     *uint          `json:"seller_id"`
	Status       string         `json:"status" gorm:"check:status IN ('in_stock','out_of_stock');not null default:'in_stock'"`
	Seller       *models.Seller `json:"seller,omitempty" gorm:"foreignKey:SellerID"`
}

type PublishBookPayload struct {
	Title        string  `json:"title" binding:"required,min=2,max=60"`
	AuthorName   string  `json:"author_name" binding:"required"`
	Genre        string  `json:"genre" binding:"required"`
	Category     string  `json:"category" binding:"required"`
	Pages        uint    `json:"pages" binding:"required,min=10"`
	Publication  string  `json:"publication" binding:"required"`
	Price        float32 `json:"price" binding:"required,min=10"`
	Units        int     `json:"units" binding:"required,min=0"`
	CoverPageUrl string  `json:"cover_page_url"`
	SellerId     *uint   `json:"seller_id"`
}

type UpdateBookPayload struct {
	Title        *string  `json:"title" binding:"omitempty,min=2,max=60"`
	AuthorName   *string  `json:"author_name"`
	Genre        *string  `json:"genre"`
	Category     *string  `json:"category"`
	Pages        *uint    `json:"pages" binding:"omitempty,min=10"`
	Publication  *string  `json:"publication"`
	Price        *float32 `json:"price" binding:"omitempty,min=10"`
	Units        *int     `json:"units" binding:"omitempty,min=0"`
	CoverPageUrl *string  `json:"cover_page_url"`
	SellerId     uint     `json:"seller_id" binding:"required"`
	BookId       uint     `json:"book_id" binding:"required"`
}

func (p *UpdateBookPayload) IsEmpty() bool {
	return p.Title == nil &&
		p.AuthorName == nil &&
		p.Genre == nil &&
		p.Category == nil &&
		p.Pages == nil &&
		p.Publication == nil &&
		p.Units == nil &&
		p.Price == nil &&
		p.CoverPageUrl == nil
}
