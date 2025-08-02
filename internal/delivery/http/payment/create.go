package payment

import (
	"encoding/json"
	"fmt"
	"github.com/adf-code/beta-payment-api/internal/delivery/response"
	"github.com/adf-code/beta-payment-api/internal/entity"
	"net/http"
)

// CreatePayment godoc
// @Summary      Create a new payment
// @Description  Creates a new payment with the status auto generate to pending
// @Tags         payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      entity.Payment  true  "Payment data to create"
// @Success      201      {object}  response.APIResponse
// @Failure      400      {object}  response.APIResponse
// @Failure      401      {object}  response.APIResponse
// @Failure      422      {object}  response.APIResponse
// @Failure      500      {object}  response.APIResponse
// @Router       /payments [post]
func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("📥 Incoming Create request")
	var payment entity.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		h.Logger.Error().Err(err).Msg("❌ Failed to store payment, invalid data")
		response.Failed(w, 422, "payments", "createPayment", "Invalid Data, Create Payment")
		return
	}

	newPayment, err := h.PaymentUC.Create(r.Context(), payment)
	if err != nil {
		h.Logger.Error().Err(err).Msg("❌ Failed to store payment, general")
		response.Failed(w, 500, "payments", "createPayment", "Error Create Payment")
		return
	}
	h.Logger.Info().Str("data", fmt.Sprint(newPayment)).Msg("✅ Successfully stored payment")
	response.Success(w, 201, "payments", "createPayment", "Success Create Payment", newPayment)
}
