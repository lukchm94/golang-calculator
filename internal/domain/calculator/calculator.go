package calculatorDomain

type OperationType string

const (
	Add       OperationType = "+"
	Substract OperationType = "-"
	Multiply  OperationType = "*"
	Divide    OperationType = "/"
)

type CalculatorInput struct {
	Number1 float64
	Number2 float64
}

type CalculationResult struct {
	Result float64 `json:"result"`
}

func Addition(c CalculatorInput) CalculationResult {
	return CalculationResult{Result: c.Number1 + c.Number2}
}

func Subtraction(c CalculatorInput) CalculationResult {
	return CalculationResult{Result: c.Number1 - c.Number2}
}

func Multiplication(c CalculatorInput) CalculationResult {
	return CalculationResult{Result: c.Number1 * c.Number2}
}

func Division(c CalculatorInput) (CalculationResult, error) {
	if c.Number2 == 0 {
		return CalculationResult{}, ErrDivisionByZero
	}
	return CalculationResult{Result: c.Number1 / c.Number2}, nil
}
