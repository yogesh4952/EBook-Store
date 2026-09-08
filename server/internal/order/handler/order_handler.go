package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/order/models"
	"github.com/yogesh4952/ebookstore/internal/order/service"
)

type OrderHandler struct {
	svc service.IOrderServ
}

func NewOrderHandler(svc service.IOrderServ) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) PlaceOrder(c *gin.Context) {

	var orderPayload models.PlaceOrderPayload

	if err := c.ShouldBindJSON(&orderPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": err.Error(),
		})
		return
	}

	userIdValue, exist := c.Get("userId")

	if exist == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Missing userid",
		})
		return
	}

	userID, ok := userIdValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "invalid user ID",
		})
		return
	}

	data, err := h.svc.PlaceOrder(c.Request.Context(), userID, orderPayload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Order placed successfully!",
		"data":    data,
	})

}
