package ssh

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/residwi/sshman/internal/interfaces"
	"github.com/residwi/sshman/utils"
)

type AgentManager struct {
	executor interfaces.CommandExecutor
}

func NewAgentManager(executor interfaces.CommandExecutor) *AgentManager {
	return &AgentManager{
		executor: executor,
	}
}

func (am *AgentManager) AddToAgent(sshPath, keyName string) error {
	keyPath := filepath.Join(sshPath, keyName)

	if utils.IsFileNotExist(keyPath) {
		return fmt.Errorf("SSH key [%s] does not exist", keyName)
	}

	if err := am.executor.Execute("ssh-add", keyPath); err != nil {
		return fmt.Errorf("failed to add key to agent: %w", err)
	}

	return nil
}

func (am *AgentManager) RemoveFromAgent(sshPath, keyName string) error {
	if isGPGAgentRunning() {
		publicKeyPath := filepath.Join(sshPath, keyName+".pub")
		if utils.IsFileNotExist(publicKeyPath) {
			return fmt.Errorf("public key [%s] does not exist", keyName)
		}

		if err := am.removeFromGPGAgent(publicKeyPath); err != nil {
			return err
		}
	} else {
		keyPath := filepath.Join(sshPath, keyName)

		if err := am.executor.Execute("ssh-add", "-d", keyPath); err != nil {
			return fmt.Errorf("failed to remove key from agent: %w", err)
		}
	}

	return nil
}

func (am *AgentManager) ListAgentKeys() ([]string, error) {
	if !am.IsAgentRunning() {
		return nil, fmt.Errorf("agent is not running")
	}

	output, err := am.executor.ExecuteWithOutput("ssh-add", "-L")
	if err != nil {
		// ssh-add -L returns exit code 1 when no keys are loaded
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list agent keys: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var keys []string
	for _, line := range lines {
		if line != "" {
			keys = append(keys, line)
		}
	}

	return keys, nil
}

func (am *AgentManager) ClearAgent(sshPath string) error {
	if isGPGAgentRunning() {
		publickeys, err := findPublicKeys(sshPath)
		if err != nil {
			return fmt.Errorf("failed to find public keys: %w", err)
		}

		for _, publicKey := range publickeys {
			if err := am.removeFromGPGAgent(publicKey); err != nil {
				return err
			}
		}
	} else {
		if err := am.executor.Execute("ssh-add", "-D"); err != nil {
			return fmt.Errorf("failed to clear agent: %w", err)
		}
	}

	return nil
}

func (am *AgentManager) IsAgentRunning() bool {
	err := am.executor.Execute("ssh-add", "-l")
	// ssh-add -l returns 2 when agent is not running
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode() != 2
	}
	return err == nil
}

func (am *AgentManager) removeFromGPGAgent(publicKeyPath string) error {
	output, err := am.executor.ExecuteWithOutput("ssh-keygen", "-lf", publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to get fingerprint of public key: %w", err)
	}

	fingerprint := strings.Split(strings.TrimSpace(string(output)), " ")[1]

	output, err = am.executor.ExecuteWithOutput("gpg-connect-agent", "keyinfo --ssh-list --ssh-fpr --with-ssh", "/bye")
	if err != nil {
		return fmt.Errorf("failed to list SSH keys in GPG agent: %w", err)
	}

	keyInfos := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")
	for keyInfo := range keyInfos {
		if strings.Contains(keyInfo, fingerprint) {
			keyGrip := strings.Split(keyInfo, " ")[2]
			deleteCmd := fmt.Sprintf("delete_key --force %s", keyGrip)

			if err := am.executor.Execute("gpg-connect-agent", deleteCmd, "/bye"); err != nil {
				return fmt.Errorf("failed to remove key from GPG agent: %w", err)
			}
		}
	}

	return nil
}

func isGPGAgentRunning() bool {
	agentPath := os.Getenv("SSH_AUTH_SOCK")
	return strings.Contains(agentPath, "gpg-agent")
}

func findPublicKeys(sshPath string) ([]string, error) {
	var publicKeys []string

	err := filepath.WalkDir(sshPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() || !strings.HasSuffix(path, ".pub") {
			return nil
		}

		publicKeys = append(publicKeys, path)

		return nil
	})

	return publicKeys, err
}
