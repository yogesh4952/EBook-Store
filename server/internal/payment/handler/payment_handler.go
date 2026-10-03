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
// ConfirmEsewaPayment godoc
// @Summary      Confirm an eSewa payment
// @Description  Browser-facing endpoint that eSewa redirects the customer to after a payment. It decodes the base64 `data` query parameter, verifies the callback signature, confirms the transaction against eSewa's status API, marks the order paid, and redirects the browser to the frontend. This endpoint is unauthenticated because the caller is the customer's browser arriving from eSewa, which cannot present a JWT.
// @Description
// @Description  On success the response is a 302 redirect to `FRONTEND_URL/payment-success?orderId=<transaction_uuid>`. On a business failure it is a 302 redirect to `FRONTEND_URL/payment-failure?reason=<reason>`, where reason is one of `invalid_signature`, `order_not_found`, `order_not_payable`, `amount_mismatch`, `payment_incomplete`, `verification_unavailable`, or `verification_failed`.
// @Description
// @Description  A 400 is returned only when the `data` parameter itself is missing or undecodable. Replaying this endpoint for an already-paid order is safe: it returns the same success redirect without re-running fulfillment.
// @Tags         payment
// @Produce      json
// @Param        data  query   string  true  "Base64-encoded JSON callback signed by eSewa"
// @Success      302   {string} string  "Redirect to FRONTEND_URL/payment-success with the order id as a query parameter"
// @Failure      400   {object} map[string]interface{}  "Missing or undecodable data parameter"
// @Router       /payment/success [get]
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
// HandleEsewaFailure godoc
// @Summary      Handle a cancelled or failed eSewa payment
// @Description  Browser-facing endpoint that eSewa redirects the customer to when they cancel at the payment page or the payment is declined. There is no payment to verify on a cancel, so the attempt is logged and the browser is redirected to the frontend. This endpoint is unauthenticated because the caller is the customer's browser arriving from eSewa.
// @Description
// @Description  The order is deliberately left in the PENDING payment status: a cancelled checkout is not a dead order, and the customer is free to retry with the same transaction uuid.
// @Tags         payment
// @Produce      json
// @Param        TransactionUuid  query  string  false  "Transaction identifier sent by eSewa (also accepted as transaction_uuid)"
// @Param        reason           query  string  false  "Failure reason reported by eSewa"
// @Success      302  {string}  string  "Redirect to FRONTEND_URL/payment-failure"
// @Router       /payment/failure [get]
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
