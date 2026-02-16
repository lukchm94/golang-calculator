package httpInfra

import (
	calculatorDomain "app/internal/domain/calculator"
	"app/internal/infrastructure/http/controllers"
	reqErr "app/internal/infrastructure/http/errors"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type CalculatorHandler struct {
	logger     *slog.Logger
	controller *controllers.CalculatorController
}

func NewCalculatorHandler(logger *slog.Logger, controller *controllers.CalculatorController) *CalculatorHandler {
	return &CalculatorHandler{logger: logger, controller: controller}
}

func (h *CalculatorHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle(string(CalculateRoute), h)
}

// ServeHTTP is a method that handles HTTP requests for the CalculatorHandler. It takes in an http.ResponseWriter and an http.Request as parameters. Inside the method, you would typically parse the request, call the appropriate methods from the calculatorService to perform calculations, and then write the response back to the client.
func (h *CalculatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received calculation request", "method", r.Method, "url", r.URL.Path)

	h.logger.Debug("Checking the Session ID form header", "session-id", r.Header.Get("X-Session-ID"), "req", r)

	ctx := r.Context()

	result, err := h.controller.Run(ctx, r)

	if err != nil {
		h.logger.Error("Calculation failed", "error", err)
		h.handleErrors(w, err)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

func (h *CalculatorHandler) handleErrors(w http.ResponseWriter, err error) {
	var invalidReqErr reqErr.InvalidRequestError
	var missingFieldError reqErr.MissingFieldError
	var invalidMethodErr reqErr.InvalidRequestMethodError

	switch {
	case errors.As(err, &missingFieldError),
		errors.As(err, &invalidReqErr),
		errors.Is(err, calculatorDomain.ErrInvalidOperation),
		errors.Is(err, calculatorDomain.ErrDivisionByZero):

		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.As(err, &invalidMethodErr):

		http.Error(w, err.Error(), http.StatusMethodNotAllowed)

	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
