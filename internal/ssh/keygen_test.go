package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/residwi/sshman/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyName(t *testing.T) {
	tests := []struct {
		name     string
		keyType  string
		purpose  string
		expected string
	}{
		{"ed25519_with_purpose", "ed25519", "work", "id_ed25519_work"},
		{"rsa_with_purpose", "rsa", "personal", "id_rsa_personal"},
		{"ed25519_without_purpose", "ed25519", "", "id_ed25519"},
		{"rsa_without_purpose", "rsa", "", "id_rsa"},
		{"rsa_with_underscore", "rsa", "my_key", "id_rsa_my_key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateKeyName(tt.keyType, tt.purpose)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestKeyGenerator_GenerateKey_ED25519_WithPurpose_Success(t *testing.T) {
	tempDir := t.TempDir()
	mockExecutor := &mocks.MockCommandExecutor{}

	expectedArgs := []string{"-t", "ed25519", "-f", filepath.Join(tempDir, "id_ed25519_work"), "-C", "test@example.com"}
	mockExecutor.On("Execute", "ssh-keygen", expectedArgs).Return(nil)

	keyGen := NewKeyGenerator(mockExecutor)

	config := KeyConfig{
		Type:    "ed25519",
		Email:   "test@example.com",
		Purpose: "work",
		SSHPath: tempDir,
	}

	keyName, err := keyGen.GenerateKey(config)

	assert.NoError(t, err)
	assert.Equal(t, "id_ed25519_work", keyName)
	mockExecutor.AssertExpectations(t)
}

func TestKeyGenerator_GenerateKey_RSA_WithoutPurpose_Success(t *testing.T) {
	tempDir := t.TempDir()
	mockExecutor := &mocks.MockCommandExecutor{}

	expectedArgs := []string{"-t", "rsa", "-f", filepath.Join(tempDir, "id_rsa"), "-b", "4096", "-C", "rsa@example.com"}
	mockExecutor.On("Execute", "ssh-keygen", expectedArgs).Return(nil)

	keyGen := NewKeyGenerator(mockExecutor)

	config := KeyConfig{
		Type:    "rsa",
		Email:   "rsa@example.com",
		Purpose: "",
		SSHPath: tempDir,
	}

	keyName, err := keyGen.GenerateKey(config)

	assert.NoError(t, err)
	assert.Equal(t, "id_rsa", keyName)
	mockExecutor.AssertExpectations(t)
}

func TestKeyGenerator_GenerateKey_KeyAlreadyExists_Error(t *testing.T) {
	tempDir := t.TempDir()
	mockExecutor := &mocks.MockCommandExecutor{}

	existingKeyPath := filepath.Join(tempDir, "id_ed25519_existing")
	err := os.WriteFile(existingKeyPath, []byte("existing key"), 0600)
	require.NoError(t, err)

	keyGen := NewKeyGenerator(mockExecutor)

	config := KeyConfig{
		Type:    "ed25519",
		Email:   "test@example.com",
		Purpose: "existing",
		SSHPath: tempDir,
	}

	_, err = keyGen.GenerateKey(config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Contains(t, err.Error(), "id_ed25519_existing")
	mockExecutor.AssertNotCalled(t, "Execute")
}

func TestKeyGenerator_GenerateKey_SSHKeygenFails_Error(t *testing.T) {
	tempDir := t.TempDir()
	mockExecutor := &mocks.MockCommandExecutor{}

	mockExecutor.On("Execute", "ssh-keygen", mock.Anything).Return(fmt.Errorf("ssh-keygen failed"))

	keyGen := NewKeyGenerator(mockExecutor)

	config := KeyConfig{
		Type:    "ed25519",
		Email:   "test@example.com",
		Purpose: "fail",
		SSHPath: tempDir,
	}

	_, err := keyGen.GenerateKey(config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create SSH key")
	assert.Contains(t, err.Error(), "ssh-keygen failed")
	mockExecutor.AssertExpectations(t)
}
