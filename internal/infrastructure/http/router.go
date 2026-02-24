package httpInfra

import (
	"net/http"
)

func NewRouter(healthHandler *HealthHandler, calculatorHandler *CalculatorHandler, userHandler *UserHandler) http.Handler {
	mux := http.NewServeMux()

	healthHandler.RegisterRoutes(mux)
	calculatorHandler.RegisterRoutes(mux)
	userHandler.RegisterRoutes(mux)

	return mux
}
