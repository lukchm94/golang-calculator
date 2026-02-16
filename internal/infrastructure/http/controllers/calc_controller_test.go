package controllers

import (
	calculatorApplication "app/internal/application/calculator"
	calculatorDomain "app/internal/domain/calculator"
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockCalculationRepository satisfies the CalculatorRepository interface
type MockCalculationRepository struct {
	// SaveFunc allows us to inject custom behavior for the Save method
	SaveFunc func(ctx context.Context, input calculatorApplication.SavedCalculationInput) error

	// Called signals if the method was actually triggered
	SaveCalled bool
}

// Save implements the interface method
func (m *MockCalculationRepository) Save(ctx context.Context, input calculatorApplication.SavedCalculationInput) error {
	m.SaveCalled = true
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, input)
	}
	return nil
}

func TestCalculatorController_Run(t *testing.T) {
	logger := slog.Default() // Uses the standard system logger
	mockRepo := &MockCalculationRepository{
		SaveFunc: func(ctx context.Context, input calculatorApplication.SavedCalculationInput) error {
			return nil // Simulate successful save
		},
	}

	service := calculatorApplication.NewCalculatorService(logger, mockRepo)

	controller := NewCalculatorController(logger, service)

	// Create a mock HTTP request
	req := httptest.NewRequest(http.MethodPost, "http://localhost/calculate", bytes.NewBufferString(`{"operation": "+", "number1": 2, "number2": 3}`))
	ctx := context.Background() // Create a background context for testing

	result, err := controller.Run(ctx, req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedValue := 5.0
	if result.Result != expectedValue {
		t.Errorf("Expected result value to be %v, got %v", expectedValue, result.Result)
	}
}

func TestCalculatorController_Run_InvalidJSON(t *testing.T) {
	logger := slog.Default() // Uses the standard system logger

	mockRepo := &MockCalculationRepository{
		SaveFunc: func(ctx context.Context, input calculatorApplication.SavedCalculationInput) error {
			return nil // Simulate successful save
		},
	}

	service := calculatorApplication.NewCalculatorService(logger, mockRepo)

	controller := NewCalculatorController(logger, service)

	// Create a mock HTTP request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "http://localhost/calculate", bytes.NewBufferString(`{"operation": "+", "number1": 2, "number2":}`))
	ctx := context.Background() // Create a background context for testing

	result, err := controller.Run(ctx, req)

	if err == nil {
		t.Fatalf("Expected error for invalid JSON, got nil")
	}

	if result != (calculatorDomain.Result{}) {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestCalculatorController_Run_InvalidMethod(t *testing.T) {
	logger := slog.Default() // Uses the standard system logger

	mockRepo := &MockCalculationRepository{
		SaveFunc: func(ctx context.Context, input calculatorApplication.SavedCalculationInput) error {
			return nil // Simulate successful save
		},
	}

	service := calculatorApplication.NewCalculatorService(logger, mockRepo)

	controller := NewCalculatorController(logger, service)

	// Create a mock HTTP request with an invalid method
	req := httptest.NewRequest(http.MethodGet, "http://localhost/calculate", nil)
	ctx := context.Background() // Create a background context for testing

	result, err := controller.Run(ctx, req)

	if err == nil {
		t.Fatalf("Expected error for invalid HTTP method, got nil")
	}

	if result != (calculatorDomain.Result{}) {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestCalculatorController_Run_InvalidOperation(t *testing.T) {
	logger := slog.Default() // Uses the standard system logger

	mockRepo := &MockCalculationRepository{
		SaveFunc: func(ctx context.Context, input calculatorApplication.SavedCalculationInput) error {
			return nil // Simulate successful save
		},
	}

	service := calculatorApplication.NewCalculatorService(logger, mockRepo)

	controller := NewCalculatorController(logger, service)

	// Create a mock HTTP request with an invalid operation
	req := httptest.NewRequest(http.MethodPost, "http://localhost/calculate", bytes.NewBufferString(`{"operation": "modulo", "number1": 10, "number2": 3}`))

	ctx := context.Background() // Create a background context for testing

	result, err := controller.Run(ctx, req)

	if err == nil {
		t.Fatalf("Expected error for invalid operation, got nil")
	}

	if result != (calculatorDomain.Result{}) {
		t.Errorf("Expected empty result, got %v", result)
	}
}

func TestCalculatorController_Run_MissingPayloadFields(t *testing.T) {
	logger := slog.Default() // Uses the standard system logger

	mockRepo := &MockCalculationRepository{
		SaveFunc: func(ctx context.Context, input calculatorApplication.SavedCalculationInput) error {
			return nil // Simulate successful save
		},
	}

	service := calculatorApplication.NewCalculatorService(logger, mockRepo)

	controller := NewCalculatorController(logger, service)

	tests := []struct {
		name string
		body string
	}{
		{
			name: "missing operation",
			body: `{"number1": 2, "number2": 3}`,
		},
		{
			name: "missing number1",
			body: `{"operation": "+", "number2": 3}`,
		},
		{
			name: "missing number2",
			body: `{"operation": "+", "number1": 2}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "http://localhost/calculate", bytes.NewBufferString(test.body))
			ctx := context.Background() // Create a background context for testing

			result, err := controller.Run(ctx, req)

			if err == nil {
				t.Fatalf("Expected error for missing payload fields, got nil")
			}

			if result != (calculatorDomain.Result{}) {
				t.Errorf("Expected empty result, got %v", result)
			}
		})
	}
}
