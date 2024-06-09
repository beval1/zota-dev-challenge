package common

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	_ "zota-dev-challenge/internal/deposit/common/zota"
	"zota-dev-challenge/internal/deposit/shared"
)

// Handler
// ClientRequest represents the request payload for a deposit.
// @Summary deposit example
// @Schemes
// @Description handle deposit
// @Tags deposit
// @Accept json
// @Produce json
// @Param depositRequest body ClientRequest true "Deposit ClientRequest"
// @Success 200 {object} _.DepositResponse "Deposit Successful"
// @Router /deposit [post]
func Handler(service ServiceInterface, logger *zap.Logger, validator *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req shared.ClientRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode request body", zap.Error(err))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		//validate request
		if err := validator.Struct(req); err != nil {
			logger.Error("Failed to validate request", zap.Error(err))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		res, err := service.ProcessDeposit(&req)
		if err != nil {
			logger.Error("Failed to process deposit", zap.Error(err))
			http.Error(w, "Failed to process deposit", http.StatusInternalServerError)
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
