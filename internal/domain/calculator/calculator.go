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

type Result struct {
	Result float64 `json:"result"`
}

func Addition(c CalculatorInput) Result {
	return Result{Result: c.Number1 + c.Number2}
}

func Subtraction(c CalculatorInput) Result {
	return Result{Result: c.Number1 - c.Number2}
}

func Multiplication(c CalculatorInput) Result {
	return Result{Result: c.Number1 * c.Number2}
}

func Division(c CalculatorInput) (Result, error) {
	if c.Number2 == 0 {
		return Result{}, ErrDivisionByZero
	}
	return Result{Result: c.Number1 / c.Number2}, nil
}
