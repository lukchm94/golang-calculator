package httpInfra

import (
	"app/internal/application"
	userService "app/internal/application/user"
	authDomain "app/internal/domain/auth"
	middleware "app/internal/infrastructure/http/middleware"
	"encoding/json"
	"log/slog"
	"net/http"
)

type HealthHandler struct {
	logger        *slog.Logger
	healthService *application.HealthService
	userService   *userService.UserService
}

func NewHealthHandler(logger *slog.Logger, healthService *application.HealthService, userService *userService.UserService) *HealthHandler {
	return &HealthHandler{logger: logger, healthService: healthService, userService: userService}
}

func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle(string(HealthRoute), h)
}

func (h *HealthHandler) RegisterAdminRoutes(adminRouter *AdminRouter) {
	adminRouter.HandleFunc(string(AdminHealthRoute), h.HandleAdminHealth)
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	healthStatus := h.healthService.CheckHealth()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthStatus)
}

func (h *HealthHandler) HandleAdminHealth(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserClaimsKey).(*authDomain.CustomClaims)

	user, err := h.userService.GetUserByID(claims.UserID)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	healthStatus := h.healthService.CheckAdminHealth(user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthStatus)
}
