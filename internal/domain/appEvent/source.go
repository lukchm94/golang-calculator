package appEvent

type EventSource string

const (
	CalculatorApp EventSource = "calculator.app"
	AuthService   EventSource = "auth.service"
)

func (s EventSource) String() string {
	return string(s)
}

func (s EventSource) IsValid() bool {
	switch s {
	case CalculatorApp, AuthService:
		return true
	}
	return false
}
