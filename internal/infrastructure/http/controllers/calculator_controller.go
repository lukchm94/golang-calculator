package controllers

import (
	calculatorApplication "app/internal/application/calculator"
	calculatorDomain "app/internal/domain/calculator"
	reqErr "app/internal/infrastructure/http/errors"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type CalculatorRequest struct {
	Operation *string  `json:"operation"`
	Number1   *float64 `json:"number1"`
	Number2   *float64 `json:"number2"`
}

type CalculatorController struct {
	logger  *slog.Logger
	service *calculatorApplication.CalculatorService
}

func NewCalculatorController(logger *slog.Logger, service *calculatorApplication.CalculatorService) *CalculatorController {
	return &CalculatorController{logger: logger, service: service}
}

func (c *CalculatorController) Run(ctx context.Context, r *http.Request) (calculatorDomain.Result, error) {
	validReq, err := c.validateReq(r)

	if err != nil {
		c.logger.Error("Request validation failed", "error", err)

		return calculatorDomain.Result{}, err
	}

	input := c.buildServiceInput(&validReq)

	c.logger.Info("Running calculation", "input", input)

	result, err := c.service.Calculate(ctx, input)

	if err == nil {
		c.logger.Info("Calculation successful", "result", result)
	}

	return result, err
}

func (c *CalculatorController) validateReq(r *http.Request) (CalculatorRequest, error) {
	c.logger.Info("Handling calculate request")

	if r.Method != http.MethodPost {
		c.logger.Error("Invalid HTTP method", "method", r.Method)

		return CalculatorRequest{}, reqErr.InvalidRequestMethodError{Method: r.Method}
	}

	req, err := c.validatePayload(r)

	if err != nil {
		c.logger.Error("Invalid request", "error", err)

		return CalculatorRequest{}, err
	}

	return *req, nil
}

func (c *CalculatorController) validatePayload(r *http.Request) (*CalculatorRequest, error) {
	var req CalculatorRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, reqErr.InvalidRequestError{Details: "invalid JSON format"}
	}

	// Ensure no fields are missing
	if req.Operation == nil {
		return nil, reqErr.MissingFieldError{FieldName: "operation"}
	}
	if req.Number1 == nil {
		return nil, reqErr.MissingFieldError{FieldName: "number1"}
	}
	if req.Number2 == nil {
		return nil, reqErr.MissingFieldError{FieldName: "number2"}
	}

	return &req, nil
}

func (c *CalculatorController) buildServiceInput(cr *CalculatorRequest) calculatorApplication.ServiceInput {
	return calculatorApplication.ServiceInput{
		Number1:  *cr.Number1,
		Number2:  *cr.Number2,
		Operator: *cr.Operation,
	}
}
