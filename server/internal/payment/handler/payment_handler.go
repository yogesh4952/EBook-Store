package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type IPaymentService interface{}

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
