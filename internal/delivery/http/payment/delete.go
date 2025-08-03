package payment

import (
	"github.com/adf-code/beta-payment-api/internal/delivery/http/router"
	"github.com/adf-code/beta-payment-api/internal/delivery/response"
	"github.com/google/uuid"
	"net/http"
)

// DeletePaymentByID godoc
// @Summary      Delete a payment by ID
// @Description  Deletes a payment entity using its UUID
// @Tags         payments
// @Security     BearerAuth
// @Param        id   path      string  true  "UUID of the payment to delete"
// @Success      202  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse  "Invalid UUID"
// @Failure      401  {object}  response.APIResponse  "Unauthorized"
// @Failure      404  {object}  response.APIResponse  "Payment not found"
// @Failure      500  {object}  response.APIResponse  "Internal server error"
// @Router       /api/v1/payments/{id} [delete]
func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("📥 Incoming Delete request")
	idStr := router.GetParam(r, "id")
	if idStr == "" {
		h.Logger.Error().Msg("❌ Failed to remove payment, missing ID parameter")
		response.Failed(w, 422, "payments", "deletePaymentByID", "Missing ID Parameter, Delete Payment by ID")
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Logger.Error().Err(err).Msg("❌ Failed to remove payment, invalid UUID parameter")
		response.Failed(w, 422, "payments", "deletePaymentByID", "Invalid UUID, Delete Payment by ID")
		return
	}
	if err := h.PaymentUC.Delete(r.Context(), id); err != nil {
		h.Logger.Error().Err(err).Msg("❌ Failed to remove payment, general")
		response.Failed(w, 500, "payments", "deletePaymentByID", "Error Delete Payment")
		return
	}
	h.Logger.Info().Msg("✅ Successfully removed payment")
	response.Success(w, 202, "payments", "deletePaymentByID", "Success Delete Payment", nil)
}
