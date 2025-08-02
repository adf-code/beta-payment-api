package payment

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/adf-code/beta-payment-api/internal/delivery/http/router"
	"github.com/adf-code/beta-payment-api/internal/delivery/response"
	"github.com/google/uuid"
	"net/http"
)

// GetPaymentByID godoc
// @Summary      Get payment by ID
// @Description  Retrieve a payment entity using its UUID
// @Tags         payments
// @Security     BearerAuth
// @Param        id   path      string  true  "UUID of the payment"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse  "Invalid UUID"
// @Failure      401  {object}  response.APIResponse  "Unauthorized"
// @Failure      404  {object}  response.APIResponse  "Payment not found"
// @Failure      500  {object}  response.APIResponse  "Internal server error"
// @Router       /payments/{id} [get]
func (h *PaymentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("📥 Incoming GetByID request")

	idStr := router.GetParam(r, "id")
	if idStr == "" {
		h.Logger.Error().Msg("❌ Failed to get payment by ID, missing ID parameter")
		response.Failed(w, 422, "payments", "getPaymentByID", "Missing ID Parameter, Get Payment by ID")
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Logger.Error().Err(err).Msg("❌ Failed to get payment by ID, invalid UUID parameter")
		response.Failed(w, 422, "payments", "getPaymentByID", "Invalid UUID, Get Payment by ID")
		return
	}
	payment, err := h.PaymentUC.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.Info().Msg("✅ Successfully get payment by id, data not found")
			response.Success(w, 404, "payments", "getPaymentByID", "Payment not Found", nil)
			return
		}
		h.Logger.Error().Err(err).Msg("❌ Failed to get payment by ID, general")
		response.Failed(w, 500, "payments", "getPaymentByID", "Error Get Payment by ID")
		return
	}
	h.Logger.Info().Str("data", fmt.Sprint(payment.ID)).Msg("✅ Successfully get payment by id")
	response.Success(w, 200, "payments", "getPaymentByID", "Success Get Payment by ID", payment)
}
