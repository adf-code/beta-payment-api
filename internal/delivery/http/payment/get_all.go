package payment

import (
	"github.com/adf-code/beta-payment-api/internal/delivery/request"
	"github.com/adf-code/beta-payment-api/internal/delivery/response"
	"net/http"
)

// GetAllPayments godoc
// @Summary      Get list of payments
// @Description  List all payments with filter, search, pagination
// @Tags         payments
// @Accept       json
// @Produce      json
//
// --- Search Query ---
// @Param        search_field      query    string   false  "Search field (e.g., title)"
// @Param        search_value      query    string   false  "Search value (e.g., golang)"
//
// --- Filter Search Query ---
// @Param filter_field query []string false "Filter field" collectionFormat(multi) explode(true)
// @Param filter_value query []string false "Filter value" collectionFormat(multi) explode(true)
//
// --- Range Query ---
// @Param range_field query []string false "Range field" collectionFormat(multi) explode(true)
// @Param from        query []string false "Range lower bound" collectionFormat(multi) explode(true)
// @Param to          query []string false "Range upper bound" collectionFormat(multi) explode(true)
//
// --- Pagination & Sort ---
// @Param        sort_field        query    string   false  "Sort field"
// @Param        sort_direction    query    string   false  "Sort direction ASC/DESC"
// @Param        page              query    int      false  "Page number"
// @Param        per_page          query    int      false  "Limit per page"
//
// @Security     BearerAuth
//
// @Success      200     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Router       /api/v1/payments [get]
func (h *PaymentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info().Msg("📥 Incoming GetAll request")
	params := request.ParsePaymentQueryParams(r)
	payments, err := h.PaymentUC.GetAll(r.Context(), params)
	if err != nil {
		h.Logger.Error().Err(err).Msg("❌ Failed to fetch payments, general")
		response.FailedWithMeta(w, 500, "payments", "getAllPayments", "Error Get All Payments", nil)
		return
	}
	h.Logger.Info().Int("count", len(payments)).Msg("✅ Successfully fetched payments")
	response.SuccessWithMeta(w, 200, "payments", "getAllPayments", "Success Get All Payments", &params, payments)
}
