package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/yogesh4952/ebookstore/internal/payment/model"
)

type IPayment interface {
	HandleSuccess(ctx context.Context) error
	HandleFailure(ctx context.Context) error
	VerifySignature(ctx context.Context, payload *model.PaymentCallback) error
}

type IPaymentRepo interface{}

type paymentService struct {
	repo IPaymentRepo
}

func NewPaymentService(repo IPaymentRepo) *paymentService {
	return &paymentService{repo: repo}
}

func (ps *paymentService) HandleSuccess(ctx context.Context) {
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

func (ps *paymentService) VerifySignature(ctx context.Context, payload *model.PaymentCallback) error {

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
		return errors.New("signature mismatch")
	}
	return nil
}
