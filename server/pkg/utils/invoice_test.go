package utils

import (
	"bytes"
	"testing"
)

func TestGenerateInvoice(t *testing.T) {
	mockOrder := OrderData{
		OrderID:      "ORD-TEST-22",
		CustomerName: "Yogesh Shah",
		TotalAmount:  2000,
	}

	pdfBytes, err := GenerateOrderInvoice(mockOrder)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("Expected PDF bytes to be generated, got empty slice")
	}
	pdfHeader := []byte("%PDF-")
	if !bytes.HasPrefix(pdfBytes, pdfHeader) {
		t.Errorf("Expected output to start with %%PDF- header, got: %s", string(pdfBytes[:5]))
	}
}
