package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/residwi/sshman/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentManager_AddToAgent_Success(t *testing.T) {
	tempDir := t.TempDir()
	mockExecutor := mocks.NewMockCommandExecutor(t)

	keyPath := filepath.Join(tempDir, "id_ed25519_test")
	err := os.WriteFile(keyPath, []byte("test key"), 0600)
	require.NoError(t, err)

	mockExecutor.EXPECT().Execute("ssh-add", []string{keyPath}).Return(nil)

	agentMgr := NewAgentManager(mockExecutor)
	err = agentMgr.AddToAgent(tempDir, "id_ed25519_test")

	assert.NoError(t, err)
}

func TestAgentManager_AddToAgent_KeyNotExists_Error(t *testing.T) {
	tempDir := t.TempDir()
	mockExecutor := mocks.NewMockCommandExecutor(t)

	agentMgr := NewAgentManager(mockExecutor)
	err := agentMgr.AddToAgent(tempDir, "nonexistent_key")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
	assert.Contains(t, err.Error(), "nonexistent_key")
}

func TestAgentManager_RemoveFromAgent_SSHAgent_Success(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/tmp/ssh-agent/socket")

	tempDir := t.TempDir()
	mockExecutor := mocks.NewMockCommandExecutor(t)

	keyPath := filepath.Join(tempDir, "id_ed25519_remove")
	os.WriteFile(keyPath+".pub", []byte("test public key"), 06400)
	mockExecutor.EXPECT().Execute("ssh-add", []string{"-d", keyPath}).Return(nil)

	agentMgr := NewAgentManager(mockExecutor)
	err := agentMgr.RemoveFromAgent(tempDir, "id_ed25519_remove")

	assert.NoError(t, err)
}

func TestAgentManager_RemoveFromAgent_GPGAgent_Success(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/tmp/gpg-agent/socket")

	tempDir := t.TempDir()
	mockExecutor := mocks.NewMockCommandExecutor(t)

	publicKeyPath := filepath.Join(tempDir, "id_rsa_gpg.pub")
	os.WriteFile(publicKeyPath, []byte("test public key for gpg"), 0644)

	fingerprintOutput := "2048 SHA256:abc123def456 test@example.com (RSA)"
	mockExecutor.EXPECT().ExecuteWithOutput("ssh-keygen", []string{"-lf", publicKeyPath}).Return([]byte(fingerprintOutput), nil)

	keyListOutput := "S KEYINFO 1234567890ABCDEF SHA256:abc123def456 - - - P - - -"
	mockExecutor.EXPECT().ExecuteWithOutput("gpg-connect-agent", []string{"keyinfo --ssh-list --ssh-fpr --with-ssh", "/bye"}).Return([]byte(keyListOutput), nil)

	mockExecutor.EXPECT().Execute("gpg-connect-agent", []string{"delete_key --force 1234567890ABCDEF", "/bye"}).Return(nil)

	agentMgr := NewAgentManager(mockExecutor)
	err := agentMgr.RemoveFromAgent(tempDir, "id_rsa_gpg")

	assert.NoError(t, err)
}

func TestAgentManager_RemoveFromAgent_GPGAgent_PublicKeyNotExists_Error(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/tmp/gpg-agent/socket")

	tempDir := t.TempDir()
	mockExecutor := mocks.NewMockCommandExecutor(t)

	agentMgr := NewAgentManager(mockExecutor)
	err := agentMgr.RemoveFromAgent(tempDir, "nonexistent_key")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "public key [nonexistent_key] does not exist")
}

func TestAgentManager_RemoveFromAgent_SSHAddFails_Error(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/tmp/ssh-agent/socket")

	tempDir := t.TempDir()
	mockExecutor := mocks.NewMockCommandExecutor(t)

	keyPath := filepath.Join(tempDir, "id_ed25519_fail")
	mockExecutor.EXPECT().Execute("ssh-add", []string{"-d", keyPath}).Return(fmt.Errorf("ssh-add failed"))

	agentMgr := NewAgentManager(mockExecutor)
	err := agentMgr.RemoveFromAgent(tempDir, "id_ed25519_fail")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ssh-add failed")
}

func TestAgentManager_ListAgentKeys_MultipleKeys_Success(t *testing.T) {
	mockExecutor := mocks.NewMockCommandExecutor(t)
	mockExecutor.EXPECT().Execute("ssh-add", []string{"-l"}).Return(nil)

	sampleOutput := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... user@example.com\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQ... user@example.com"
	mockExecutor.EXPECT().ExecuteWithOutput("ssh-add", []string{"-L"}).Return([]byte(sampleOutput), nil)

	agentMgr := NewAgentManager(mockExecutor)
	keys, err := agentMgr.ListAgentKeys()

	assert.NoError(t, err)
	assert.Len(t, keys, 2)
	assert.Contains(t, keys[0], "ssh-ed25519")
	assert.Contains(t, keys[1], "ssh-rsa")
}

func TestAgentManager_ListAgentKeys_NoKeys_Success(t *testing.T) {
	mockExecutor := mocks.NewMockCommandExecutor(t)

	mockExecutor.EXPECT().Execute("ssh-add", []string{"-l"}).Return(nil)
	mockExecutor.EXPECT().ExecuteWithOutput("ssh-add", []string{"-L"}).Return([]byte(""), nil)

	agentMgr := NewAgentManager(mockExecutor)
	keys, err := agentMgr.ListAgentKeys()

	assert.NoError(t, err)
	assert.Empty(t, keys)
}

func TestAgentManager_ListAgentKeys_AgentNotRunning_Error(t *testing.T) {
	mockExecutor := mocks.NewMockCommandExecutor(t)

	mockExecutor.EXPECT().Execute("ssh-add", []string{"-l"}).Return(fmt.Errorf("ssh-add: agent is not running"))

	agentMgr := NewAgentManager(mockExecutor)
	keys, err := agentMgr.ListAgentKeys()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "agent is not running")
	assert.Nil(t, keys)
}

func TestAgentManager_ListAgentKeys_SSHAddFails_Error(t *testing.T) {
	mockExecutor := mocks.NewMockCommandExecutor(t)

	mockExecutor.EXPECT().Execute("ssh-add", []string{"-l"}).Return(nil)
	mockExecutor.EXPECT().ExecuteWithOutput("ssh-add", []string{"-L"}).Return([]byte(""), fmt.Errorf("unexpected error"))

	agentMgr := NewAgentManager(mockExecutor)
	keys, err := agentMgr.ListAgentKeys()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list agent keys")
	assert.Contains(t, err.Error(), "unexpected error")
	assert.Nil(t, keys)
}

func TestAgentManager_ClearAgent_SSHAgent_Success(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/tmp/ssh-agent/socket")

	tempDir := t.TempDir()
	mockExecutor := mocks.NewMockCommandExecutor(t)

	mockExecutor.EXPECT().Execute("ssh-add", []string{"-D"}).Return(nil)

	agentMgr := NewAgentManager(mockExecutor)
	err := agentMgr.ClearAgent(tempDir)

	assert.NoError(t, err)
}

func TestAgentManager_ClearAgent_GPGAgent_Success(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/tmp/gpg-agent/socket")

	tempDir := t.TempDir()
	mockExecutor := mocks.NewMockCommandExecutor(t)

	publicKey := filepath.Join(tempDir, "id_ed25519.pub")
	os.WriteFile(publicKey, []byte("test public key for gpg"), 0644)

	mockExecutor.EXPECT().ExecuteWithOutput("ssh-keygen", []string{"-lf", publicKey}).Return([]byte("256 SHA256:abc123 test@example.com (ED25519)"), nil)
	mockExecutor.EXPECT().ExecuteWithOutput("gpg-connect-agent", []string{"keyinfo --ssh-list --ssh-fpr --with-ssh", "/bye"}).Return([]byte("S KEYINFO ABCDEF123456 SHA256:abc123 - - - P - - -"), nil)
	mockExecutor.EXPECT().Execute("gpg-connect-agent", []string{"delete_key --force ABCDEF123456", "/bye"}).Return(nil)

	agentMgr := NewAgentManager(mockExecutor)
	err := agentMgr.ClearAgent(tempDir)

	assert.NoError(t, err)
}

func TestAgentManager_ClearAgent_Fails_Error(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/tmp/ssh-agent/socket")

	tempDir := t.TempDir()
	mockExecutor := mocks.NewMockCommandExecutor(t)

	mockExecutor.EXPECT().Execute("ssh-add", []string{"-D"}).Return(fmt.Errorf("ssh-add failed"))

	agentMgr := NewAgentManager(mockExecutor)
	err := agentMgr.ClearAgent(tempDir)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to clear agent")
	assert.Contains(t, err.Error(), "ssh-add failed")
}

func TestAgentManager_IsAgentRunning(t *testing.T) {
	mockExecutor := mocks.NewMockCommandExecutor(t)

	mockExecutor.EXPECT().Execute("ssh-add", []string{"-l"}).Return(nil)

	agentMgr := NewAgentManager(mockExecutor)
	isRunning := agentMgr.IsAgentRunning()

	assert.True(t, isRunning)
}

func TestIsGPGAgentRunning(t *testing.T) {
	tests := []struct {
		name     string
		authSock string
		expected bool
	}{
		{"gpg_agent_running", "/tmp/gpg-agent/socket", true},
		{"gpg_agent_with_path", "/home/user/.gnupg/S.gpg-agent.ssh", true},
		{"ssh_agent_running", "/tmp/ssh-agent/socket", false},
		{"openssh_agent", "/tmp/ssh-XXXXXXXX/agent.12345", false},
		{"empty_auth_sock", "", false},
		{"no_auth_sock", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SSH_AUTH_SOCK", tt.authSock)

			result := isGPGAgentRunning()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindPublicKeys(t *testing.T) {
	tempDir := t.TempDir()

	testFiles := []struct {
		name     string
		isPublic bool
	}{
		{"id_rsa", false},
		{"id_rsa.pub", true},
		{"id_ed25519", false},
		{"id_ed25519.pub", true},
		{"config", false},
		{"known_hosts", false},
		{"id_ecdsa.pub", true},
		{"not_a_key.txt", false},
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file.name)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		require.NoError(t, err)
	}

	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(subDir, "id_rsa_backup.pub"), []byte("backup key"), 0644)
	assert.NoError(t, err)

	publicKeys, err := findPublicKeys(tempDir)
	require.NoError(t, err)

	assert.Len(t, publicKeys, 4)

	expectedKeys := []string{
		filepath.Join(tempDir, "id_rsa.pub"),
		filepath.Join(tempDir, "id_ed25519.pub"),
		filepath.Join(tempDir, "id_ecdsa.pub"),
		filepath.Join(subDir, "id_rsa_backup.pub"),
	}

	for _, expectedKey := range expectedKeys {
		assert.Contains(t, publicKeys, expectedKey)
	}
}

func TestFindPublicKeys_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	publicKeys, err := findPublicKeys(tempDir)

	assert.NoError(t, err)
	assert.Empty(t, publicKeys)
}

func TestFindPublicKeys_NonexistentDirectory(t *testing.T) {
	nonexistentDir := "/nonexistent/directory"

	publicKeys, err := findPublicKeys(nonexistentDir)

	assert.Error(t, err)
	assert.Nil(t, publicKeys)
}
