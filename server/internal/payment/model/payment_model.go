package model

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// PaymentCallback is the payload eSewa base64-encodes and sends to the
// success_url query string as ?data=
type PaymentCallback struct {
	TransactionCode  string `json:"transaction_code"`
	Status           string `json:"status"`
	TotalAmount      string `json:"total_amount"`
	TransactionUUID  string `json:"transaction_uuid"`
	ProductCode      string `json:"product_code"`
	SignedFieldNames string `json:"signed_field_names"`
	Signature        string `json:"signature"`
}

// EsewaStatusResponse is the JSON body returned by eSewa's server-to-server
// transaction status check endpoint.
type EsewaStatusResponse struct {
	ProductCode      string     `json:"product_code"`
	TransactionCode  string     `json:"transaction_code"`
	TransactionUUID  string     `json:"transaction_uuid"`
	TotalAmount      FlexAmount `json:"total_amount"`
	Status           string     `json:"status"`
	RefId            string     `json:"ref_id"`
	SignedFieldNames string     `json:"signed_field_names"`
	Signature        string     `json:"signature"`
}

// FlexAmount tolerates the two shapes eSewa uses for money. The status-check
// endpoint answers with a JSON number (1750.0) while the checkout form and the
// browser callback both use a JSON string ("1750.0"). Decoding into a plain
// string fails outright on the number, so accept either.
type FlexAmount string

func (f *FlexAmount) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)

	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*f = ""
		return nil
	}

	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return err
		}
		*f = FlexAmount(s)
		return nil
	}

	// Numeric form: re-render without an exponent so it stays comparable to the
	// signed string form of the same value.
	value, err := strconv.ParseFloat(string(trimmed), 64)
	if err != nil {
		return err
	}
	*f = FlexAmount(strconv.FormatFloat(value, 'f', -1, 64))

	return nil
}

func (f FlexAmount) String() string {
	return string(f)
}

// EsewaProductCode is the merchant product code registered with eSewa.
const EsewaProductCode = "EPAYTEST"

// EsewaStatusComplete is the terminal success state reported by eSewa.
const EsewaStatusComplete = "COMPLETE"
