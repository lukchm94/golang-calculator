package httpInfra

import (
	userDomain "app/internal/domain/user"
	"app/internal/infrastructure/http/controllers"
	reqErr "app/internal/infrastructure/http/errors"
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
	mux.HandleFunc(string(RegisterRoute), h.handleRegister)
	mux.HandleFunc(string(LoginRoute), h.handleLogin)
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logInvalidMethod(r)

		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	ctx := r.Context()

	res, err := h.controller.Login(ctx, r)

	if err != nil {
		h.logger.Error("User Login failed", "error", err)
		h.handleErrors(w, err)

		return
	}

	h.logger.Info("Logged in user", "user", res)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(userDomain.User{
		ID:        res.ID,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
	})
}

func (h *UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logInvalidMethod(r)

		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info("Received user request", "method", r.Method, "url", r.URL.Path)

	ctx := r.Context()

	result, err := h.controller.Register(ctx, r)

	if err != nil {
		h.logger.Error("User registration failed", "error", err)
		h.handleErrors(w, err)

		return
	}

	h.setHeaderContentTypeApplicationJson(w)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(userDomain.User{
		ID:        result.ID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	})
}

func (h *UserHandler) handleErrors(w http.ResponseWriter, err error) {
	h.setHeaderContentTypeApplicationJson(w)

	status := http.StatusInternalServerError

	switch err.(type) {
	case reqErr.InvalidRequestError, reqErr.MissingFieldError:
		// Invalid or incomplete request payload
		status = http.StatusBadRequest
	case reqErr.InvalidRequestMethodError:
		// Wrong HTTP method used for this endpoint
		status = http.StatusMethodNotAllowed
	case reqErr.NotImplementedError:
		// Functionality (e.g. login) not implemented yet
		status = http.StatusNotImplemented
	case reqErr.UserNotFoundError:
		status = http.StatusNotFound

	case reqErr.InvalidCredentialsError:
		status = http.StatusUnauthorized
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func (h *UserHandler) logInvalidMethod(r *http.Request) {
	h.logger.Error("Invalid HTTP method. Required POST and received: ", "requestMethod", r.Method)
}

func (h *UserHandler) setHeaderContentTypeApplicationJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
