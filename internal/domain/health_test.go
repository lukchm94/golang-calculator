package domain

import (
	"testing"
)

func TestHealthCheck(t *testing.T) {
	result := HealthCheck()

	if result.Status != HealthyStatus {
		t.Errorf("Expected status %q, got %q", HealthyStatus, result.Status)
	}
}

func TestHealthStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		status   HealthStatusResponse
		expected string
	}{
		{
			name:     "HealthyStatus constant",
			status:   HealthyStatus,
			expected: "healthy",
		},
		{
			name:     "ErrorStatus constant",
			status:   ErrorStatus,
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, tt.status)
			}
		})
	}
}
