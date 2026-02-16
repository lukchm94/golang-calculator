package calculatorDomain

import (
	"testing"
)

func TestAddition(t *testing.T) {
	input := CalculatorInput{Number1: 5, Number2: 3}
	result := Addition(input)

	expected := 8.0
	if result.Result != expected {
		t.Errorf("Expected %f, got %f", expected, result.Result)
	}
}

func TestSubtraction(t *testing.T) {
	input := CalculatorInput{Number1: 5, Number2: 3}
	result := Subtraction(input)

	expected := 2.0
	if result.Result != expected {
		t.Errorf("Expected %f, got %f", expected, result.Result)
	}
}

func TestMultiplication(t *testing.T) {
	input := CalculatorInput{Number1: 5, Number2: 3}
	result := Multiplication(input)

	expected := 15.0
	if result.Result != expected {
		t.Errorf("Expected %f, got %f", expected, result.Result)
	}
}

func TestDivision(t *testing.T) {
	input := CalculatorInput{Number1: 5, Number2: 2}
	result, err := Division(input)

	expected := 2.5
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.Result != expected {
		t.Errorf("Expected %f, got %f", expected, result.Result)
	}
}

func TestDivisionByZero(t *testing.T) {
	input := CalculatorInput{Number1: 5, Number2: 0}
	_, err := Division(input)

	if err == nil {
		t.Error("Expected an error for division by zero, but got none")
	}
}

func TestOperationTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		opType   OperationType
		expected string
	}{
		{
			name:     "Add constant",
			opType:   Add,
			expected: "+",
		},
		{
			name:     "Substract constant",
			opType:   Substract,
			expected: "-",
		},
		{
			name:     "Multiply constant",
			opType:   Multiply,
			expected: "*",
		},
		{
			name:     "Divide constant",
			opType:   Divide,
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.opType) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, tt.opType)
			}
		})
	}
}
