package calculatorApplication

import (
	"context"
	"log/slog"
	"testing"
)

// MockCalculationRepository satisfies the CalculatorRepository interface
type MockCalculationRepository struct {
	// SaveFunc allows us to inject custom behavior for the Save method
	SaveFunc func(ctx context.Context, input SavedCalculationInput) error

	// Called signals if the method was actually triggered
	SaveCalled bool
}

// Save implements the interface method
func (m *MockCalculationRepository) Save(ctx context.Context, input SavedCalculationInput) error {
	m.SaveCalled = true
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, input)
	}
	return nil
}
func TestCalculatorService_Calculate(t *testing.T) {
	logger := slog.Default()

	tests := []struct {
		name          string
		input         ServiceInput
		expectedValue float64
		expectError   bool
		mockErr       error
	}{
		{
			name: "Addition",
			input: ServiceInput{
				Number1:  2,
				Number2:  3,
				Operator: "+",
			},
			expectedValue: 5,
			expectError:   false,
		},
		{
			name: "Subtraction",
			input: ServiceInput{
				Number1:  5,
				Number2:  3,
				Operator: "-",
			},
			expectedValue: 2,
			expectError:   false,
		},
		{
			name: "Multiplication",
			input: ServiceInput{
				Number1:  4,
				Number2:  3,
				Operator: "*",
			},
			expectedValue: 12,
			expectError:   false,
		},
		{
			name: "Division",
			input: ServiceInput{
				Number1:  10,
				Number2:  2,
				Operator: "/",
			},
			expectedValue: 5,
			expectError:   false,
		},
		{
			name: "Division by Zero",
			input: ServiceInput{
				Number1:  10,
				Number2:  0,
				Operator: "/",
			},
			expectedValue: 0,
			expectError:   true,
		},

		{
			name: "Invalid Operation",
			input: ServiceInput{
				Number1:  10,
				Number2:  2,
				Operator: "modulo",
			},
			expectedValue: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Initialize the mock for this specific case
			mockRepo := &MockCalculationRepository{
				SaveFunc: func(ctx context.Context, input SavedCalculationInput) error {
					return tt.mockErr
				},
			}

			service := NewCalculatorService(logger, mockRepo)
			result, err := service.Calculate(ctx, tt.input)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if !tt.expectError && result.Result != tt.expectedValue {
				t.Errorf("Expected result %v, got %v", tt.expectedValue, result.Result)
			}

			// Peer Tip: Add a check to see if SaveCalculation was called on success!
			if !tt.expectError && !mockRepo.SaveCalled {
				t.Errorf("Expected SaveCalculation to be called, but it wasn't")
			}
		})
	}
}
