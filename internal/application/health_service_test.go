package application

import (
	"testing"
)

func TestHealthService(t *testing.T) {
	healthService := NewHealthService()
	healthStatus := healthService.CheckHealth()

	if healthStatus.Status != "healthy" {
		t.Errorf("Expected health status to be 'healthy', got '%s'", healthStatus.Status)
	}
}
