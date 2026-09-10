package handler

import (
	"errors"
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

// PlaceOrder godoc
// @Summary      Place a new order
// @Description  Create a new order with items, payment method, and shipping address
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        payload body models.PlaceOrderPayload true "Order data"
// @Success      202 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /order/place-order [post]
// @Security     BearerAuth
func (h *OrderHandler) PlaceOrder(c *gin.Context) {

	var orderPayload models.PlaceOrderPayload

	if err := c.ShouldBindJSON(&orderPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": "false",
			"message": err.Error(),
		})
		return
	}

	userIdValue, _ := c.Get("userId")

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
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, models.ErrInvalidAddress), errors.Is(err, models.ErrBookNotFound):
			status = http.StatusNotFound
		case errors.Is(err, models.ErrInvalidPaymentMethod):
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Order placed succesfully",
		"data":    data,
	})

}

func (h *OrderHandler) ListUserOrder(c *gin.Context) {
	userId, _ := c.Get("userId")

	userIdValue, ok := userId.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "invalid user ID",
		})
		return
	}

	data, err := h.svc.ListUserOrder(c.Request.Context(), userIdValue)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "Order fetch successfully",
		"data":    data,
	})

}
