package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/cart/service"
)

type CartHandler struct {
	cartServ service.ICartService
}

func NewCartHandler(serv service.ICartService) *CartHandler {
	return &CartHandler{cartServ: serv}
}

func (ch *CartHandler) AddToCart(ctx *gin.Context) {}
