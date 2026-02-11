package config

import (
	"testing"
)

func TestDefaultPort(t *testing.T) {
	if DefaultPort != "8080" {
		t.Errorf("Expected DefaultPort to be '8080', got '%s'", DefaultPort)
	}
}

func TestEmptyString(t *testing.T) {
	if EmptyString != "" {
		t.Errorf("Expected EmptyString to be '', got '%s'", EmptyString)
	}
}
