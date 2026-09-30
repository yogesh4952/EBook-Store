package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/yogesh4952/ebookstore/internal/payment/model"
)

func TestVerifyWithEsewa(t *testing.T) {
	tests := []struct {
		name           string
		gatewayStatus  int
		gatewayBody    string
		expectedAmount float32
		wantErr        error
	}{
		{
			// eSewa's status endpoint answers with a JSON *number*, unlike the
			// checkout form and the browser callback which use strings.
			name:           "complete and amount matches as json number",
			gatewayStatus:  http.StatusOK,
			gatewayBody:    `{"product_code":"EPAYTEST","transaction_uuid":"ORD-J7NYJZBDSA","total_amount":1750.0,"status":"COMPLETE","ref_id":"000H8NP"}`,
			expectedAmount: 1750,
			wantErr:        nil,
		},
		{
			name:           "complete and amount matches as json string",
			gatewayStatus:  http.StatusOK,
			gatewayBody:    `{"product_code":"EPAYTEST","transaction_uuid":"ORD-J7NYJZBDSA","total_amount":"1750.0","status":"COMPLETE","ref_id":"000H8NP"}`,
			expectedAmount: 1750,
			wantErr:        nil,
		},
		{
			name:           "complete but amount differs",
			gatewayStatus:  http.StatusOK,
			gatewayBody:    `{"total_amount":99.0,"status":"COMPLETE"}`,
			expectedAmount: 1750,
			wantErr:        model.ErrAmountMismatch,
		},
		{
			name:           "status not complete",
			gatewayStatus:  http.StatusOK,
			gatewayBody:    `{"total_amount":1750.0,"status":"PENDING"}`,
			expectedAmount: 1750,
			wantErr:        model.ErrStatusNotComplete,
		},
		{
			name:           "gateway http error",
			gatewayStatus:  http.StatusInternalServerError,
			expectedAmount: 1750,
			wantErr:        model.ErrEsewaUnavailable,
		},
		{
			name:           "unparseable amount",
			gatewayStatus:  http.StatusOK,
			gatewayBody:    `{"total_amount":"not-a-number","status":"COMPLETE"}`,
			expectedAmount: 1750,
			wantErr:        model.ErrAmountMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotQuery url.Values

			gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotQuery = r.URL.Query()
				w.WriteHeader(tt.gatewayStatus)
				if tt.gatewayStatus == http.StatusOK {
					_, _ = w.Write([]byte(tt.gatewayBody))
				}
			}))
			defer gateway.Close()

			t.Setenv("ESEWA_BASE_URL", gateway.URL)

			serv := NewPaymentService(&fakeRepo{})
			err := serv.VerifyWithEsewa(context.Background(), &model.PaymentCallback{
				ProductCode:     model.EsewaProductCode,
				TotalAmount:     "1750.0",
				TransactionUUID: "ORD-J7NYJZBDSA",
			}, tt.expectedAmount)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got %v, want %v", err, tt.wantErr)
			}

			// The amount must reach eSewa byte-for-byte as it was signed,
			// never reformatted by %f into "1750.000000".
			if got := gotQuery.Get("total_amount"); got != "1750.0" {
				t.Errorf("total_amount sent as %q, want %q", got, "1750.0")
			}
			if got := gotQuery.Get("transaction_uuid"); got != "ORD-J7NYJZBDSA" {
				t.Errorf("transaction_uuid sent as %q", got)
			}
		})
	}
}

func TestVerifyWithEsewa_MissingBaseURL(t *testing.T) {
	t.Setenv("ESEWA_BASE_URL", "")

	serv := NewPaymentService(&fakeRepo{})
	err := serv.VerifyWithEsewa(context.Background(), &model.PaymentCallback{}, 1750)

	if !errors.Is(err, model.ErrEsewaUnavailable) {
		t.Fatalf("expected ErrEsewaUnavailable, got %v", err)
	}
}
