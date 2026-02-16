package calculatorApplication

import (
	"log/slog"
	"testing"
)

func TestCalculatorService_Calculate(t *testing.T) {
	logger := slog.Default() // Uses the standard system logger

	service := NewCalculatorService(logger)

	tests := []struct {
		name          string
		input         ServiceInput
		expectedValue float64
		expectError   bool
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
			result, err := service.Calculate(tt.input)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !tt.expectError {
				if result.Result != tt.expectedValue {
					t.Errorf("Expected result value to be %v, got %v", tt.expectedValue, result.Result)
				}
			}
		})
	}
}
