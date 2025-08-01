package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSupportedKeyType(t *testing.T) {
	tests := []struct {
		name     string
		keyType  string
		expected bool
	}{
		{"Valid RSA key", "rsa", true},
		{"Valid ed25519 key", "ed25519", true},
		{"invalid key type", "invalid", false},
		{"Empty key type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			error := isSupportedKeyType(tt.keyType)
			if tt.expected {
				assert.NoError(t, error, "Expected no error for key type %s", tt.keyType)
			} else {
				assert.Error(t, error, "Expected error for key type %s", tt.keyType)
			}
		})
	}
}
