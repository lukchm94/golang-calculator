package httpInfra

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/application"
	domain "app/internal/domain/health"
)

func TestHealthHandler(t *testing.T) {
	logger := slog.Default() // Uses the standard system logger

	// Create a new instance of the HealthService
	healthService := application.NewHealthService(logger)
	healthHandler := NewHealthHandler(logger, healthService)

	// Create a mock HTTP request
	req := httptest.NewRequest(http.MethodGet, "http://localhost/health", nil)

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call ServeHTTP
	healthHandler.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the content type
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	// Decode the response body
	var response domain.HealthStatus
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	// Assert that the response status is healthy
	if response.Status != domain.HealthyStatus {
		t.Errorf("Expected health status %q, got %q", domain.HealthyStatus, response.Status)
	}
}
