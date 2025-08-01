package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProviderConfig(t *testing.T) {
	tests := []struct {
		name         string
		provider     string
		expectedUser string
		expectedHost string
		shouldExist  bool
	}{
		{"github_valid", "github", "git", "github.com", true},
		{"gitlab_valid", "gitlab", "git", "gitlab.com", true},
		{"bitbucket_valid", "bitbucket", "git", "bitbucket.org", true},
		{"invalid_provider", "invalid", "", "", false},
		{"empty_provider", "", "", "", false},
		{"generic_provider", "generic", "", "", false},
		{"uppercase_github", "GITHUB", "", "", false},
		{"mixed_case_gitlab", "GitLab", "", "", false},
		{"nonexistent_provider", "nonexistent", "", "", false},
		{"numeric_provider", "123", "", "", false},
		{"special_chars_provider", "git@hub", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, exists := GetProviderConfig(tt.provider)

			assert.Equal(t, tt.shouldExist, exists)

			if tt.shouldExist {
				assert.Equal(t, tt.expectedUser, config.User)
				assert.Equal(t, tt.expectedHost, config.Hostname)
			} else {
				assert.Equal(t, "", config.User)
				assert.Equal(t, "", config.Hostname)
			}
		})
	}
}

func TestGetSupportedProviders(t *testing.T) {
	providers := GetSupportedProviders()

	expectedProviders := []string{"github", "gitlab", "bitbucket", "generic"}
	assert.Equal(t, len(expectedProviders), len(providers))

	for _, expected := range expectedProviders {
		assert.Contains(t, providers, expected)
	}

	assert.Equal(t, expectedProviders, providers)
}
