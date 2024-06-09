package common

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"net/http"
	"zota-dev-challenge/internal/status/shared"

	_ "zota-dev-challenge/internal/status/common/zota"
)

var decoder = schema.NewDecoder()

// Handler
// ClientRequest represents the request payload for a deposit.
// @Summary status check example
// @Schemes
// @Description handle status check
// @Tags status check
// @Accept json
// @Produce json
// @Param orderId query string true "Order ID"
// @Param merchantOrderId query string true "Merchant Order ID"
// @Success 200 {object} _.StatusResponse "Status Check successful"
// @Router /status [get]
func Handler(service ServiceInterface, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req shared.ClientRequest
		// Decode request params - it's get request
		if err := decoder.Decode(&req, r.URL.Query()); err != nil {
			logger.Error("Failed to decode request", zap.Error(err))
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
			return
		}

		res, err := service.CheckStatus(&req)
		if err != nil {
			logger.Error("Failed to check status", zap.Error(err))
			http.Error(w, "Failed to check status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			logger.Error("Failed to encode response", zap.Error(err))
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
