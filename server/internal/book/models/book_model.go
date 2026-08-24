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

//
