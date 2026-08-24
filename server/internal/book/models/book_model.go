package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title        string  `json:"title"`
	AuthorName   string  `json:"author_name"`
	Genre        string  `json:"genre"`
	Category     string  `json:"category"`
	Pages        uint    `json:"pages"`
	Publication  string  `json:"publication"`
	Price        float32 `json:"price"`
	Units        int     `json:"units"`
	CoverPageUrl string  `json:"cover_page_url"`
	SellerID     *uint   `json:"seller_id"`
}

type BookPayload struct {
	Title        string  `json:"title" binding:"required"`
	AuthorName   string  `json:"author_name" binding:"required"`
	Genre        string  `json:"genre" binding:"required"`
	Category     string  `json:"category" binding:"required"`
	Pages        uint    `json:"pages" binding:"required"`
	Publication  string  `json:"publication" binding:"required"`
	Price        float32 `json:"price" binding:"required"`
	Units        int     `json:"units" binding:"required"`
	CoverPageUrl string  `json:"cover_page_url" `
}
