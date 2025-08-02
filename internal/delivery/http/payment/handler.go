package payment

import (
	"github.com/adf-code/beta-payment-api/internal/usecase"
	"github.com/rs/zerolog"
)

type PaymentHandler struct {
	PaymentUC usecase.PaymentUseCase
	Logger    zerolog.Logger
}

func NewPaymentHandler(paymentUC usecase.PaymentUseCase, logger zerolog.Logger) *PaymentHandler {
	return &PaymentHandler{PaymentUC: paymentUC, Logger: logger}
}
