package httpInfra

import (
	"encoding/json"
	"net/http"

	"app/internal/application"
)

type HealthHandler struct {
	healthService *application.HealthService
}

func NewHealthHandler(healthService *application.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	healthStatus := h.healthService.CheckHealth()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthStatus)
}
