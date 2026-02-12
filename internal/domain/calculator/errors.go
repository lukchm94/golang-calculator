package calculatorDomain

import "errors"

type CalculatorError error

var (
	ErrDivisionByZero   CalculatorError = errors.New("division by zero is not allowed")
	ErrInvalidOperation CalculatorError = errors.New("invalid operation type")
)
