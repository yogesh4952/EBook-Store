package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"testing"

	orderModel "github.com/yogesh4952/ebookstore/internal/order/models"
	"github.com/yogesh4952/ebookstore/internal/payment/model"
)

// sign mirrors what utils.GenerateSignature produces, so the test exercises the
// same message construction the real checkout flow uses.
func sign(t *testing.T, p *model.PaymentCallback) string {
	t.Helper()

	parts := make([]string, 0)
	for _, name := range strings.Split(p.SignedFieldNames, ",") {
		value, ok := fieldValue(p, strings.TrimSpace(name))
		if !ok {
			t.Fatalf("unknown signed field %q", name)
		}
		parts = append(parts, strings.TrimSpace(name)+"="+value)
	}

	mac := hmac.New(sha256.New, []byte(os.Getenv("ESEWA_SECRET")))
	mac.Write([]byte(strings.Join(parts, ",")))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

type fakeRepo struct {
	order     *orderModel.Order
	findErr   error
	findCalls int
	confirmed bool
	txCode    string
}

func (f *fakeRepo) FindOrderByTransactionUUID(_ context.Context, _ string) (*orderModel.Order, error) {
	f.findCalls++
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.order, nil
}

func (f *fakeRepo) ConfirmPayment(_ context.Context, _ uint, transactionCode string) error {
	f.confirmed = true
	f.txCode = transactionCode
	return nil
}

func paidOrder() *orderModel.Order {
	return &orderModel.Order{
		PaymentMethod: orderModel.PayementEsewa,
		PaymentStatus: orderModel.PayementPaid,
		TotalPrice:    1750,
	}
}

func pendingOrder() *orderModel.Order {
	return &orderModel.Order{
		PaymentMethod: orderModel.PayementEsewa,
		PaymentStatus: orderModel.PayementPending,
		TotalPrice:    1750,
	}
}

// signedCallback returns a payload with a valid local signature, so tests that
// are not about signature checking get past step 2.
func signedCallback(t *testing.T) *model.PaymentCallback {
	t.Helper()

	t.Setenv("ESEWA_SECRET", "s3cr3t")

	payload := &model.PaymentCallback{
		TransactionCode:  "000H8NP",
		Status:           model.EsewaStatusComplete,
		TotalAmount:      "1750.0",
		TransactionUUID:  "ORD-J7NYJZBDSA",
		ProductCode:      model.EsewaProductCode,
		SignedFieldNames: "transaction_code,status,total_amount,transaction_uuid,product_code,signed_field_names",
	}
	payload.Signature = sign(t, payload)

	return payload
}

func TestConfirmEsewaPayment_SkipsEsewaWhenAlreadyPaid(t *testing.T) {
	repo := &fakeRepo{order: paidOrder()}
	serv := NewPaymentService(repo)
	// No gateway is configured: if the already-paid path leaked through to
	// verification, this would error instead of returning nil.
	t.Setenv("ESEWA_BASE_URL", "")

	err := serv.ConfirmEsewaPayment(context.Background(), signedCallback(t))

	if err != nil {
		t.Fatalf("expected idempotent success, got %v", err)
	}
	if repo.confirmed {
		t.Error("already-paid order must not be confirmed again")
	}
}

func TestConfirmEsewaPayment_StopsOnBadSignature(t *testing.T) {
	t.Setenv("ESEWA_SECRET", "test-secret")

	repo := &fakeRepo{order: pendingOrder()}
	serv := NewPaymentService(repo)

	err := serv.ConfirmEsewaPayment(context.Background(), &model.PaymentCallback{
		TransactionUUID:  "ORD-1",
		SignedFieldNames: "total_amount,transaction_uuid,product_code",
		Signature:        "forged",
		TotalAmount:      "1750.0",
	})

	if !errors.Is(err, model.ErrSignatureMismatch) {
		t.Fatalf("expected ErrSignatureMismatch, got %v", err)
	}
	if repo.findCalls != 0 {
		t.Error("must reject a bad signature before touching the database")
	}
}

func TestConfirmEsewaPayment_PropagatesOrderNotFound(t *testing.T) {
	repo := &fakeRepo{findErr: model.ErrOrderNotFound}
	serv := NewPaymentService(repo)

	err := serv.ConfirmEsewaPayment(context.Background(), signedCallback(t))

	if !errors.Is(err, model.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
	if repo.confirmed {
		t.Error("must not confirm a missing order")
	}
}

func TestConfirmEsewaPayment_RejectsCashOnDelivery(t *testing.T) {
	repo := &fakeRepo{order: &orderModel.Order{
		PaymentMethod: orderModel.PayementCOD,
		PaymentStatus: orderModel.PayementPending,
		TotalPrice:    1750,
	}}
	serv := NewPaymentService(repo)

	err := serv.ConfirmEsewaPayment(context.Background(), signedCallback(t))

	if !errors.Is(err, model.ErrOrderNotPayable) {
		t.Fatalf("expected ErrOrderNotPayable, got %v", err)
	}
}

func TestVerifyEsewaSignature_AcceptsValidSignature(t *testing.T) {
	t.Setenv("ESEWA_SECRET", "s3cr3t")

	serv := NewPaymentService(&fakeRepo{})
	payload := &model.PaymentCallback{
		TransactionUUID:  "ORD-J7NYJZBDSA",
		TotalAmount:      "1750.0",
		ProductCode:      model.EsewaProductCode,
		SignedFieldNames: "total_amount,transaction_uuid,product_code",
	}
	// Signed with the same helper checkout uses.
	payload.Signature = sign(t, payload)

	if err := serv.verifyEsewaSignature(payload); err != nil {
		t.Fatalf("expected valid signature to pass, got %v", err)
	}
}

func TestVerifyEsewaSignature_RejectsUnknownField(t *testing.T) {
	t.Setenv("ESEWA_SECRET", "s3cr3t")

	serv := NewPaymentService(&fakeRepo{})
	payload := &model.PaymentCallback{
		SignedFieldNames: "total_amount,attacker_field",
		Signature:        "x",
	}

	if err := serv.verifyEsewaSignature(payload); err == nil {
		t.Fatal("expected unknown signed field to be rejected")
	}
}
