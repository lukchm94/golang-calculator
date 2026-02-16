package calculatorApplication

import (
	calculatorDomain "app/internal/domain/calculator"
	"context"
)

type SavedCalculationInput struct {
	Operation        calculatorDomain.OperationType
	CalculationInput calculatorDomain.CalculatorInput
	Result           calculatorDomain.Result
	SessionId        string
}

type CalculationRepository interface {
	Save(ctx context.Context, input SavedCalculationInput) error
}
