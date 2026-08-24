package auth

import "errors"

var (
	ErrDuplicateEmail = errors.New("Email address is already registerd")
	ErrInvalid        = errors.New("Inalid role; must be cutomer, admin, vendor")
)
