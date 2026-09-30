package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yogesh4952/ebookstore/internal/payment/model"
	"github.com/yogesh4952/ebookstore/pkg/logger"
)

type IPaymentService interface {
	ConfirmEsewaPayment(ctx context.Context, payload *model.PaymentCallback) error
}

type PaymentHandler struct {
	paymentServ IPaymentService
	frontendURL string
}

func NewPaymentHandler(paymentServ IPaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentServ: paymentServ,
		frontendURL: os.Getenv("FRONTEND_URL"),
	}
}

// ConfirmEsewaPayment is the browser-facing endpoint eSewa redirects to after a
// payment. It only translates: decode the callback, hand it to the service, then
// bounce the user to the SPA. Malformed input is a 400; business failures redirect.
func (h *PaymentHandler) ConfirmEsewaPayment(c *gin.Context) {
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

	if err := h.paymentServ.ConfirmEsewaPayment(c.Request.Context(), &payload); err != nil {
		reason := failureReason(err)

		logger.Ctx(c.Request.Context()).Error().Err(err).
			Str("transaction_uuid", payload.TransactionUUID).
			Str("reason", reason).
			Msg("esewa payment confirmation failed")

		c.Redirect(http.StatusFound, h.failureURL(reason))
		return
	}

	c.Redirect(http.StatusFound, h.frontendURL+"/payment-success?orderId="+url.QueryEscape(payload.TransactionUUID))
}

// HandleEsewaFailure is where eSewa sends the browser when the customer
// cancels or the payment is declined. There is nothing to verify on a cancel,
// so this only logs the attempt and bounces the user back to the SPA.
//
// The order is deliberately left PENDING: a cancelled checkout is not a dead
// order, and the customer is free to retry with the same transaction uuid.
func (h *PaymentHandler) HandleEsewaFailure(c *gin.Context) {
	transactionUUID := c.Query("TransactionUuid")
	if transactionUUID == "" {
		transactionUUID = c.Query("transaction_uuid")
	}

	logger.Ctx(c.Request.Context()).Warn().
		Str("transaction_uuid", transactionUUID).
		Str("reason", c.Query("reason")).
		Msg("esewa payment failed or was cancelled")

	c.Redirect(http.StatusFound, h.frontendURL+"/payment-failure")
}

func (h *PaymentHandler) failureURL(reason string) string {
	return h.frontendURL + "/payment-failure?reason=" + url.QueryEscape(reason)
}

func failureReason(err error) string {
	switch {
	case errors.Is(err, model.ErrSignatureMismatch):
		return "invalid_signature"
	case errors.Is(err, model.ErrOrderNotFound):
		return "order_not_found"
	case errors.Is(err, model.ErrOrderNotPayable):
		return "order_not_payable"
	case errors.Is(err, model.ErrAmountMismatch):
		return "amount_mismatch"
	case errors.Is(err, model.ErrStatusNotComplete):
		return "payment_incomplete"
	case errors.Is(err, model.ErrEsewaUnavailable):
		return "verification_unavailable"
	default:
		return "verification_failed"
	}
}
