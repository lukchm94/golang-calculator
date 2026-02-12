package httpInfra

import (
	"net/http"
)

func NewRouter(healthHandler *HealthHandler, calculatorHandler *CalculatorHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(string(HealthRoute), healthHandler)
	mux.Handle(string(CalculateRoute), calculatorHandler)

	return mux
}
