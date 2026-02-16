package calculatorApplication

import (
	calculatorDomain "app/internal/domain/calculator"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
)

type ServiceInput struct {
	Number1  float64
	Number2  float64
	Operator string
}

type CalculatorService struct {
	logger *slog.Logger
	repo   CalculationRepository
}

func NewCalculatorService(logger *slog.Logger, repo CalculationRepository) *CalculatorService {
	return &CalculatorService{logger: logger, repo: repo}
}

func (s *CalculatorService) Calculate(ctx context.Context, input ServiceInput) (calculatorDomain.Result, error) {
	s.logger.Info("Performing calculation", "number1", input.Number1, "number2", input.Number2, "operation", input.Operator)

	domainInput, operation, err := s.convertToDomainInput(input)

	if err != nil {
		return calculatorDomain.Result{}, err
	}

	result, err := s.getResult(domainInput, operation)
	if err != nil {
		return result, err
	}

	s.logger.Info("Calculation performed successfully", "result", result.Result)
	s.logger.Debug("Saving record to database")

	inputToSave := s.generateInput(operation, domainInput, result)

	err = s.repo.Save(ctx, inputToSave)

	if err != nil {
		s.logger.Error("Failed to save calculation record", "error", err)
		return result, err
	}

	s.logger.Info("Record saved successfully")
	return result, nil
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

func (s *CalculatorService) getResult(input calculatorDomain.CalculatorInput, operation calculatorDomain.OperationType) (calculatorDomain.Result, error) {
	switch operation {

	case calculatorDomain.Add:
		return calculatorDomain.Addition(input), nil

	case calculatorDomain.Substract:
		return calculatorDomain.Subtraction(input), nil

	case calculatorDomain.Multiply:
		return calculatorDomain.Multiplication(input), nil

	case calculatorDomain.Divide:
		return calculatorDomain.Division(input)

	default:
		s.logger.Error("Invalid operation type", "operation", operation)

		return calculatorDomain.Result{}, calculatorDomain.ErrInvalidOperation
	}
}

func (s *CalculatorService) generateInput(operation calculatorDomain.OperationType, calcInput calculatorDomain.CalculatorInput, result calculatorDomain.Result) SavedCalculationInput {
	s.logger.Debug("Converting the input to hex code", "input", calcInput, "operation", operation, "result", result)

	data := fmt.Sprintf("%v-%v-%v-%v", operation, calcInput.Number1, calcInput.Number2, result.Result)

	hash := sha256.Sum256([]byte(data))

	return SavedCalculationInput{
		Operation:        operation,
		CalculationInput: calcInput,
		Result:           result,
		SessionId:        hex.EncodeToString(hash[:]),
	}

}
