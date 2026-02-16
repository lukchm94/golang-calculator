package application

import (
	"log/slog"
	"testing"
)

func TestHealthService(t *testing.T) {
	logger := slog.Default() // Uses the standard system logger

	healthService := NewHealthService(logger)
	healthStatus := healthService.CheckHealth()

	if healthStatus.Status != "healthy" {
		t.Errorf("Expected health status to be 'healthy', got '%s'", healthStatus.Status)
	}
}
