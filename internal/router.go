package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	deposit "zota-dev-challenge/internal/deposit/common"
	status "zota-dev-challenge/internal/status/common"
)

func InitRouterV1(depositService *deposit.Service, statusService *status.Service, validator *validator.Validate, logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/api/v1/deposit", deposit.Handler(depositService, logger, validator))
	r.Get("/api/v1/status", status.Handler(statusService, logger))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	return r
}
