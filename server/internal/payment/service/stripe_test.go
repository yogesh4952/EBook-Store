package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	webhookUrl    = "http://localhost:8080/webhook/stripe"
	webhookSecret = "terimakichud12345"
)

func TestStripeWebhook(t *testing.T) {
	payload := []byte(`{
		"id":"evt_test_123",
		"object": "event",
		"type": "payment_intent.succeeded",
		"data": {
			"object": "pi_test_456",
			"amount": 1000,
			"currency": "usd",
			"status": "succeeded"
		}
	}, 
	"created": 1234567890
}`)

	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	fmt.Printf("Payload: %s\n", string(payload))
	fmt.Printf("Signature: %s\n", signature)

	req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(payload))

	if err != nil {

		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Singature-256", signature)

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)

	}
	defer resp.Body.Close()
	t.Logf("Response Status: %s\n", resp.Status)

}
