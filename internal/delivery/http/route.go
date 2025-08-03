package http

import (
	"github.com/adf-code/beta-payment-api/internal/delivery/http/health"
	"github.com/adf-code/beta-payment-api/internal/delivery/http/middleware"
	"github.com/adf-code/beta-payment-api/internal/delivery/http/payment"
	"github.com/adf-code/beta-payment-api/internal/delivery/http/router"
	"github.com/adf-code/beta-payment-api/internal/usecase"
	"github.com/rs/zerolog"

	"github.com/swaggo/http-swagger"
	"net/http"
)

func SetupHandler(paymentUC usecase.PaymentUseCase, logger zerolog.Logger) http.Handler {
	paymentHandler := payment.NewPaymentHandler(paymentUC, logger)
	healthHandler := health.NewHealthHandler(logger)
	auth := middleware.AuthMiddleware(logger)
	log := middleware.LoggingMiddleware(logger)

	r := router.NewRouter()

	r.HandlePrefix(http.MethodGet, "/swagger/", httpSwagger.WrapHandler)

	r.Handle("GET", "/healthz", middleware.Chain(log)(healthHandler.Check))

	r.Handle("PUT", "/api/v1/payments/status/{id}", middleware.Chain(log, auth)(paymentHandler.UpdateByID))
	r.Handle("GET", "/api/v1/payments/{id}", middleware.Chain(log, auth)(paymentHandler.GetByID))
	r.Handle("GET", "/api/v1/payments", middleware.Chain(log, auth)(paymentHandler.GetAll))
	r.Handle("POST", "/api/v1/payments", middleware.Chain(log, auth)(paymentHandler.Create))
	r.Handle("DELETE", "/api/v1/payments/{id}", middleware.Chain(log, auth)(paymentHandler.Delete))

	return r
}
