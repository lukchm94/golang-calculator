package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"app/cmd/config"
	"app/internal/application"
	calculatorApplication "app/internal/application/calculator"
	httpInfra "app/internal/infrastructure/http"
	"app/internal/infrastructure/http/controllers"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Services initialization
	healthService := application.NewHealthService()
	calculatorService := calculatorApplication.NewCalculatorService(logger)

	// Controllers initialization
	calculatorController := controllers.NewCalculatorController(logger)

	// HTTP Handlers initialization
	healthHandler := httpInfra.NewHealthHandler(healthService)
	calculatorHandler := httpInfra.NewCalculatorHandler(logger, calculatorController, calculatorService)
	// Router initialization
	router := httpInfra.NewRouter(healthHandler, calculatorHandler)

	port := os.Getenv(string(config.PortEnvKey))
	if port == string(config.EmptyString) {
		port = string(config.DefaultPort)
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Find the health endpoint at http://localhost:%s/health", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
