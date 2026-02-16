package httpInfra

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"app/internal/application"
)

type HealthHandler struct {
	logger        *slog.Logger
	healthService *application.HealthService
}

func NewHealthHandler(logger *slog.Logger, healthService *application.HealthService) *HealthHandler {
	return &HealthHandler{logger: logger, healthService: healthService}
}

func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle(string(HealthRoute), h)
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	healthStatus := h.healthService.CheckHealth()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthStatus)
}
