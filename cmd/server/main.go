package main

import (
	"log"
	"net/http"
	"os"

	"app/cmd/config"
	"app/internal/application"
	httpInfra "app/internal/infrastructure/http"
)

func main() {
	// Services initialization
	healthService := application.NewHealthService()

	// HTTP Handlers initialization
	healthHandler := httpInfra.NewHealthHandler(healthService)

	// Router initialization
	router := httpInfra.NewRouter(healthHandler)

	port := os.Getenv(string(config.PortEnvKey))
	if port == string(config.EmptyString) {
		port = string(config.DefaultPort)
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Find the health endpoint at http://localhost:%s/health", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
