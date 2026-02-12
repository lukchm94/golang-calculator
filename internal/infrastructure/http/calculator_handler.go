package httpInfra

import (
	calculatorApplication "app/internal/application/calculator"
	"app/internal/infrastructure/http/controllers"
	"encoding/json"
	"log/slog"
	"net/http"
)

type CalculatorHandler struct {
	logger     *slog.Logger
	controller *controllers.CalculatorController
	service    *calculatorApplication.CalculatorService
}

func NewCalculatorHandler(logger *slog.Logger, controller *controllers.CalculatorController, service *calculatorApplication.CalculatorService) *CalculatorHandler {
	return &CalculatorHandler{logger: logger, controller: controller, service: service}
}

// ServeHTTP is a method that handles HTTP requests for the CalculatorHandler. It takes in an http.ResponseWriter and an http.Request as parameters. Inside the method, you would typically parse the request, call the appropriate methods from the calculatorService to perform calculations, and then write the response back to the client.
func (h *CalculatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received calculation request", "method", r.Method, "url", r.URL.Path)

	validatedRequest, err := h.controller.ValidateCalculatorRequest(w, r)

	if err != nil {
		h.logger.Error("Request validation failed", "error", err)
		return
	}

	input := h.buildServiceInput(&validatedRequest)

	h.logger.Info("Validated request", "number1", input.Number1, "number2", input.Number2, "operation", input.Operator)

	result, err := h.service.Calculate(input)

	if err != nil {
		h.logger.Error("Calculation failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *CalculatorHandler) buildServiceInput(v *controllers.CalculatorRequest) calculatorApplication.ServiceInput {
	return calculatorApplication.ServiceInput{
		Number1:  *v.Number1,
		Number2:  *v.Number2,
		Operator: *v.Operation,
	}
}
