package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type CalculatorRequest struct {
	Operation *string  `json:"operation"`
	Number1   *float64 `json:"number1"`
	Number2   *float64 `json:"number2"`
}

type CalculatorController struct {
	logger *slog.Logger
}

func NewCalculatorController(logger *slog.Logger) *CalculatorController {
	return &CalculatorController{logger: logger}
}

func (c *CalculatorController) ValidateCalculatorRequest(w http.ResponseWriter, r *http.Request) (CalculatorRequest, error) {
	c.logger.Info("Handling calculate request")

	if r.Method != http.MethodPost {
		c.logger.Error("Invalid HTTP method", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return CalculatorRequest{}, http.ErrNotSupported
	}

	req, err := c.validateRequest(r)

	if err != nil {
		c.logger.Error("Invalid request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return CalculatorRequest{}, err
	}

	return *req, nil
}

func (c *CalculatorController) validateRequest(r *http.Request) (*CalculatorRequest, error) {
	var req CalculatorRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.New("invalid JSON format")
	}

	// Ensure no fields are missing
	if req.Operation == nil {
		return nil, errors.New("missing field: operation")
	}
	if req.Number1 == nil {
		return nil, errors.New("missing field: number1")
	}
	if req.Number2 == nil {
		return nil, errors.New("missing field: number2")
	}

	return &req, nil
}
