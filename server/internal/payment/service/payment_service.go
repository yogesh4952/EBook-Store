package service

import "context"

type IPayment interface {
	VetifySignature(ctx context.Context) error
	HandleSuccess(ctx context.Context) error
	HandleFailure(ctx context.Context) error
}

type IPaymentRepo interface{}

type paymentService struct {
	repo IPaymentRepo
}

func NewPaymentService(repo IPaymentRepo) *paymentService {
	return &paymentService{repo: repo}
}

// http://localhost:8080/api/payment/success
//?data=eyJ0cmFuc2FjdGlvbl9jb2RlIjoiMDAwSDNBQyIsInN0YXR1cyI6IkNPTVBMRVRFIiwidG90YWxfYW1vdW50IjoiMTc1MC4wIiwidHJhbnNhY3Rpb25fdXVpZCI6Ik9SRC1KUlBHV1NNS1VBIiwicHJvZHVjdF9jb2RlIjoiRVBBWVRFU1QiLCJzaWduZWRfZmllbGRfbmFtZXMiOiJ0cmFuc2FjdGlvbl9jb2RlLHN0YXR1cyx0b3RhbF9hbW91bnQsdHJhbnNhY3Rpb25fdXVpZCxwcm9kdWN0X2NvZGUsc2lnbmVkX2ZpZWxkX25hbWVzIiwic2lnbmF0dXJlIjoidmpJMlhlbm9lWlZ3KzJVRlAyWE5nUHpmQXA0UjJJaE5wRDdhK3MxSjJ6QT0ifQ==

func (ps *paymentService) HandleSuccess(ctx context.Context) {
}
