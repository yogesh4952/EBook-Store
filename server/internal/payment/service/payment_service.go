package service

import "context"

type IPayment interface {
	VetifySignature(ctx context.Context) error
}

type IPaymentRepo interface{}

type paymentService struct {
	repo IPaymentRepo
}

func NewPaymentService(repo IPaymentRepo) *paymentService {
	return &paymentService{repo: repo}
}

func (serv *paymentService) VerifySingature(ctx context.Context) error {

	// boyd:= ctx.
	// return nil
}
