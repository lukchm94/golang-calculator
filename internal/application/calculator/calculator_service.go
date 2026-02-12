package calculatorApplication

import (
	calculatorDomain "app/internal/domain/calculator"
	"log/slog"
)

type ServiceInput struct {
	Number1  float64
	Number2  float64
	Operator string
}

type CalculatorService struct {
	logger *slog.Logger
}

func NewCalculatorService(logger *slog.Logger) *CalculatorService {
	return &CalculatorService{logger: logger}
}

func (s *CalculatorService) Calculate(input ServiceInput) (calculatorDomain.CalculationResult, error) {
	s.logger.Info("Performing calculation", "number1", input.Number1, "number2", input.Number2, "operation", input.Operator)

	domainInput, operation, err := s.convertToDomainInput(input)

	if err != nil {
		return calculatorDomain.CalculationResult{}, err
	}

	switch operation {

	case calculatorDomain.Add:
		return calculatorDomain.Addition(domainInput), nil

	case calculatorDomain.Substract:
		return calculatorDomain.Subtraction(domainInput), nil

	case calculatorDomain.Multiply:
		return calculatorDomain.Multiplication(domainInput), nil

	case calculatorDomain.Divide:
		return calculatorDomain.Division(domainInput)

	default:
		s.logger.Error("Invalid operation type", "operation", input.Operator)

		return calculatorDomain.CalculationResult{}, calculatorDomain.ErrInvalidOperation
	}
}

func (s *CalculatorService) convertToDomainInput(input ServiceInput) (calculatorDomain.CalculatorInput, calculatorDomain.OperationType, error) {
	domainInput := calculatorDomain.CalculatorInput{
		Number1: input.Number1,
		Number2: input.Number2,
	}

	var operation calculatorDomain.OperationType

	switch input.Operator {
	case "+":
		operation = calculatorDomain.Add
	case "-":
		operation = calculatorDomain.Substract
	case "*":
		operation = calculatorDomain.Multiply
	case "/":
		operation = calculatorDomain.Divide
	default:
		s.logger.Error("Invalid operation type", "operation", input.Operator)
		return domainInput, "", calculatorDomain.ErrInvalidOperation
	}

	return domainInput, operation, nil
}
