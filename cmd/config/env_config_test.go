package config

import (
	"testing"
)

func TestPortEnvKey(t *testing.T) {
	if PortEnvKey != "PORT" {
		t.Errorf("Expected PortEnvKey to be 'PORT', got '%s'", PortEnvKey)
	}
}
