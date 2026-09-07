package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/order/service"
)

type OrderHandler struct {
	svc service.IOrderServ
}

func NewBookHandler(svc service.IOrderServ) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) PlaceOrder(c *gin.Context) {

}
