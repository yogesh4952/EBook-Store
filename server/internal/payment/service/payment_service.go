package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	orderModel "github.com/yogesh4952/ebookstore/internal/order/models"
	"github.com/yogesh4952/ebookstore/internal/payment/model"
	"github.com/yogesh4952/ebookstore/pkg/logger"
)

const (
	esewaTimeout    = 5 * time.Second
	esewaAmountTol  = 0.01
	esewaStatusPath = "/api/epay/transaction/status/"
)

type IPaymentRepo interface {
	FindOrderByTransactionUUID(ctx context.Context, transactionUUID string) (*orderModel.Order, error)
	ConfirmPayment(ctx context.Context, orderID uint, transactionCode string) error
}

type paymentService struct {
	repo   IPaymentRepo
	client *http.Client
}

func NewPaymentService(repo IPaymentRepo) *paymentService {
	return &paymentService{
		repo:   repo,
		client: &http.Client{Timeout: esewaTimeout},
	}
}

func fieldValue(p *model.PaymentCallback, name string) (string, bool) {
	switch name {
	case "transaction_code":
		return p.TransactionCode, true
	case "status":
		return p.Status, true
	case "total_amount":
		return p.TotalAmount, true
	case "transaction_uuid":
		return p.TransactionUUID, true
	case "product_code":
		return p.ProductCode, true
	case "signed_field_names":
		return p.SignedFieldNames, true
	default:
		return "", false
	}
}

func (ps *paymentService) verifyEsewaSignature(payload *model.PaymentCallback) error {
	fields := strings.Split(payload.SignedFieldNames, ",")

	parts := make([]string, 0, len(fields))
	for _, raw := range fields {
		name := strings.TrimSpace(raw)
		value, ok := fieldValue(payload, name)
		if !ok {
			return fmt.Errorf("unknown signed field %q", name)
		}
		parts = append(parts, name+"="+value)
	}
	message := strings.Join(parts, ",")
	secret := os.Getenv("ESEWA_SECRET")
	if secret == "" {
		return errors.New("ESEWA_SECRET not set")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	generated := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(generated), []byte(payload.Signature)) {
		return model.ErrSignatureMismatch
	}
	return nil
}

// ConfirmEsewaPayment runs the full callback flow in the one order that keeps
// it safe: verify the local signature, load the order, short-circuit if already
// paid, ask eSewa directly whether the money actually moved, then persist.
func (ps *paymentService) ConfirmEsewaPayment(ctx context.Context, payload *model.PaymentCallback) error {
	if err := ps.verifyEsewaSignature(payload); err != nil {
		return err
	}

	order, err := ps.repo.FindOrderByTransactionUUID(ctx, payload.TransactionUUID)
	if err != nil {
		return err
	}

	// Idempotency: a refresh or a repeated redirect must not re-run fulfillment.
	if order.PaymentStatus == orderModel.PayementPaid {
		return nil
	}

	if order.PaymentMethod != orderModel.PayementEsewa {
		return fmt.Errorf("%w: order %d is %s", model.ErrOrderNotPayable, order.ID, order.PaymentMethod)
	}

	if err := ps.VerifyWithEsewa(ctx, payload, order.TotalPrice); err != nil {
		return err
	}

	return ps.repo.ConfirmPayment(ctx, order.ID, payload.TransactionCode)
}

// VerifyWithEsewa asks eSewa's own API whether the transaction completed. The
// local signature only proves the payload was not tampered with in transit; only
// eSewa can confirm the money moved. expectedAmount comes from our database,
// never from the callback, so the comparison is not circular.
func (ps *paymentService) VerifyWithEsewa(ctx context.Context, p *model.PaymentCallback, expectedAmount float32) error {
	base := strings.TrimRight(os.Getenv("ESEWA_BASE_URL"), "/")
	if base == "" {
		return fmt.Errorf("%w: ESEWA_BASE_URL not set", model.ErrEsewaUnavailable)
	}

	// Send the amount exactly as eSewa signed it; reformatting with %f would
	// turn "1750.0" into "1750.000000" and the gateway would not match it.
	query := url.Values{}
	query.Set("product_code", p.ProductCode)
	query.Set("total_amount", p.TotalAmount)
	query.Set("transaction_uuid", p.TransactionUUID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+esewaStatusPath, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", model.ErrEsewaUnavailable, err)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := ps.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", model.ErrEsewaUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status %d", model.ErrEsewaUnavailable, resp.StatusCode)
	}

	var result model.EsewaStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("%w: decoding gateway response: %v", model.ErrEsewaUnavailable, err)
	}

	if result.Status != model.EsewaStatusComplete {
		return fmt.Errorf("%w: %s", model.ErrStatusNotComplete, result.Status)
	}

	remote, err := strconv.ParseFloat(result.TotalAmount.String(), 64)
	if err != nil {
		return fmt.Errorf("%w: unparseable %q", model.ErrAmountMismatch, result.TotalAmount.String())
	}
	if math.Abs(remote-float64(expectedAmount)) > esewaAmountTol {
		return fmt.Errorf("%w: gateway %v order %v", model.ErrAmountMismatch, remote, expectedAmount)
	}

	logger.Ctx(ctx).Info().
		Str("transaction_uuid", p.TransactionUUID).
		Str("transaction_code", result.TransactionCode).
		Msg("esewa transaction verified")

	return nil
}
