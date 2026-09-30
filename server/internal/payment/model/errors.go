package model

import "errors"

var (
	ErrSignatureMismatch = errors.New("signature mismatch")
	ErrOrderNotFound     = errors.New("order not found")
	ErrOrderNotPayable   = errors.New("order is not payable via esewa")
	ErrAmountMismatch    = errors.New("amount mismatch")
	ErrStatusNotComplete = errors.New("esewa status not complete")
	ErrEsewaUnavailable  = errors.New("esewa verification unavailable")
)
