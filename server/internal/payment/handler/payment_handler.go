package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/payment/model"
)

type IPaymentService interface {
	VerifySignature(ctx context.Context, payload *model.PaymentCallback) error
}

type PaymentHandler struct {
	paymentServ IPaymentService
}

func NewPaymentHandler(serv IPaymentService) *PaymentHandler {
	return &PaymentHandler{paymentServ: serv}
}

func (h *PaymentHandler) HandleRedirect(c *gin.Context) {
	data := c.Query("data")
	fmt.Printf("%s", data)
}

func (h *PaymentHandler) HandleSuccess(c *gin.Context) {
	data := c.Query("data")
	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "missing data"})
		return
	}

	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var payload model.PaymentCallback
	if err := json.Unmarshal(raw, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// verify signature BEFORE trusting any of this
	if err := h.paymentServ.VerifySignature(c.Request.Context(), &payload); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}
func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	data := c.Query("data")
	fmt.Printf("%s", data)
}
