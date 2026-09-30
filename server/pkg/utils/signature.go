package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func GenerateSignature(totalPrice, TransactionUUID, productCode string) (string, error) {
	message := fmt.Sprintf("total_amount=%s,transaction_uuid=%s,product_code=%s", totalPrice, TransactionUUID, productCode)
	log.Printf("%s", message)
	secretKey := os.Getenv("ESEWA_SECRET")

	if secretKey == "" {
		return "", fmt.Errorf("empty secret key")
	}

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	rawSignature := mac.Sum(nil)

	signature := base64.StdEncoding.EncodeToString(rawSignature)
	return signature, nil
}
