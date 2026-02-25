package httpInfra

import (
	userDomain "app/internal/domain/user"
	"app/internal/infrastructure/http/controllers"
	"encoding/json"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	logger     *slog.Logger
	controller *controllers.UserController
}

func NewUserHandler(logger *slog.Logger, controller *controllers.UserController) *UserHandler {
	return &UserHandler{logger: logger, controller: controller}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle(string(RegisterRoute), h)
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received user request", "method", r.Method, "url", r.URL.Path)

	ctx := r.Context()

	result, err := h.controller.Register(ctx, r)

	if err != nil {
		h.logger.Error("User registration failed", "error", err)
		h.handleErrors(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(userDomain.User{
		ID:        result.ID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	})
}

func (h *UserHandler) handleErrors(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
