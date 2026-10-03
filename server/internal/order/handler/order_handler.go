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
// @Description  Create a new order with items, payment method, and shipping address. For the `COD` payment method the order is created as PENDING. For the `ESEWA` method the response additionally carries a signed `esewa_payload` that the client must POST to eSewa's payment form; the order stays PENDING until the eSewa callback is confirmed at `/payment/success`. Prices are always recalculated server-side from the book records, so any price in the request is ignored.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        payload body models.PlaceOrderPayload true "Order data"
// @Success      202 {object} map[string]interface{}  "Order placed, containing the order under data. For ESEWA orders this also includes the signed esewa_payload"
// @Failure      400 {object} map[string]interface{}  "Invalid request body, or an unsupported payment method"
// @Failure      401 {object} map[string]interface{}  "Missing or invalid token"
// @Failure      404 {object} map[string]interface{}  "Address id or one of the book ids was not found"
// @Failure      500 {object} map[string]interface{}  "User id on the request context had an unexpected type, or the order could not be created"
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

// ListUserOrder godoc
// @Summary      List User order
// @Description  List every order belonging to the authenticated user, including their line items. The authenticated user id is taken from the JWT and cannot be overridden in the request.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      202 {object} map[string]interface{}  "Orders fetched, containing the list of orders under data"
// @Failure      400 {object} map[string]interface{}  "Orders could not be fetched"
// @Failure      401 {object} map[string]interface{}  "Missing or invalid token"
// @Failure      500 {object} map[string]interface{}  "User id on the request context had an unexpected type"
// @Router       /order/list-user-order [get]
// @Security     BearerAuth
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
