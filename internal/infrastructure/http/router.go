package httpInfra

import (
	"net/http"
)

func NewRouter(healthHandler *HealthHandler, calculatorHandler *CalculatorHandler) http.Handler {
	mux := http.NewServeMux()

	healthHandler.RegisterRoutes(mux)
	calculatorHandler.RegisterRoutes(mux)

	return mux
}
