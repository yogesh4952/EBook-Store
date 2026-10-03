package model

import "github.com/yogesh4952/ebookstore/internal/book/models"

type Cart struct {
	Id         uint `json:"id"`
	UserId     uint `json:"user_id"`
	TotalPrice int64

	Items []CartItem
}

type AddToCartPayload struct {
	BookId   uint `json:"book_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required"`
}

type ItemResp struct {
	Data     models.Book `json:"data"`
	Quantity uint        `json:"quantity"`
	Subtotal float32     `json:"subtotal"`
}
type CartResponse struct {
	Items []ItemResp `json:"items"`
	Total float32    `json:"total"`
}
