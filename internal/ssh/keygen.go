package ssh

import (
	"fmt"
	"path/filepath"

	"github.com/residwi/sshman/internal/interfaces"
	"github.com/residwi/sshman/utils"
)

type KeyConfig struct {
	Type     string
	Email    string
	Purpose  string
	Provider string
	SSHPath  string
}

type KeyGenerator struct {
	executor interfaces.CommandExecutor
}

func NewKeyGenerator(executor interfaces.CommandExecutor) *KeyGenerator {
	return &KeyGenerator{
		executor: executor,
	}
}

func (kg *KeyGenerator) GenerateKey(config KeyConfig) (string, error) {
	keyName := generateKeyName(config.Type, config.Purpose)
	filePath := filepath.Join(config.SSHPath, keyName)

	if !utils.IsFileNotExist(filePath) {
		return "", fmt.Errorf("SSH key [%s] already exists. Please choose a different purpose or delete the existing key", keyName)
	}

	var keygenArgs []string
	keygenArgs = append(keygenArgs, "-t", config.Type, "-f", filePath)

	if config.Type == "rsa" {
		keygenArgs = append(keygenArgs, "-b", "4096")
	}

	keygenArgs = append(keygenArgs, "-C", config.Email)

	if err := kg.executor.Execute("ssh-keygen", keygenArgs...); err != nil {
		return "", fmt.Errorf("failed to create SSH key: %w", err)
	}

	return keyName, nil
}

func generateKeyName(keyType, purpose string) string {
	keyName := "id_" + keyType
	if purpose != "" {
		keyName += "_" + purpose
	}
	return keyName
}
