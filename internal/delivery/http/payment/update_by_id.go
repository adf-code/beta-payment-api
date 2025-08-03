package payment

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/adf-code/beta-payment-api/internal/delivery/http/router"
	"github.com/adf-code/beta-payment-api/internal/delivery/request"
	"github.com/adf-code/beta-payment-api/internal/delivery/response"
	"github.com/google/uuid"
	"net/http"
)

// UpdatePaymentByID godoc
// @Summary      Update payment by ID
// @Description  Modify a payment entity using its UUID
// @Tags         payments
// @Security     BearerAuth
// @Param        id   path      string  true  "UUID of the payment"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse  "Invalid UUID"
// @Failure      401  {object}  response.APIResponse  "Unauthorized"
// @Failure      404  {object}  response.APIResponse  "Payment not found"
// @Failure      500  {object}  response.APIResponse  "Internal server error"
// @Router       /api/v1/payments/{id} [get]
func (h *PaymentHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("📥 Incoming UpdateByID request")

	// Step 1: Get ID from path
	idStr := router.GetParam(r, "id")
	if idStr == "" {
		h.Logger.Error().Msg("❌ Missing ID parameter")
		response.Failed(w, 422, "payments", "updatePaymentByID", "Missing ID Parameter")
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Logger.Error().Err(err).Msg("❌ Invalid UUID parameter")
		response.Failed(w, 422, "payments", "updatePaymentByID", "Invalid UUID")
		return
	}

	var req request.UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error().Err(err).Msg("❌ Failed to decode request body")
		response.Failed(w, 400, "payments", "updatePaymentByID", "Invalid Request Body")
		return
	}

	if err := req.Validate(); err != nil {
		h.Logger.Error().Err(err).Msg("❌ Validation error")
		response.Failed(w, 422, "payments", "updatePaymentByID", "Validation Error")
		return
	}

	updatedPayment, err := h.PaymentUC.UpdateByID(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.Info().Msg("✅ Payment not found for update")
			response.Success(w, 404, "payments", "updatePaymentByID", "Payment not Found", nil)
			return
		}
		h.Logger.Error().Err(err).Msg("❌ Update failed")
		response.Failed(w, 500, "payments", "updatePaymentByID", "Failed to Update Payment")
		return
	}

	h.Logger.Info().Str("id", id.String()).Msg("✅ Successfully updated payment")
	response.Success(w, 200, "payments", "updatePaymentByID", "Payment Updated", updatedPayment)
}
