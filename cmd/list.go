package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/residwi/sshman/internal/interfaces"
	"github.com/residwi/sshman/internal/ssh"
	"github.com/residwi/sshman/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List SSH keys and configurations",
	Long:  `Display all SSH keys in the SSH directory and their status`,
	RunE:  listSSHKeys,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listSSHKeys(cmd *cobra.Command, args []string) error {
	sshPath, _ := cmd.Flags().GetString("ssh-path")

	privateKeys, err := findPrivateKeys(sshPath)
	if err != nil {
		return fmt.Errorf("failed to list SSH keys: %w", err)
	}

	if len(privateKeys) == 0 {
		utils.PrintSuccess("No SSH keys found in " + sshPath)
		return nil
	}

	executor := &interfaces.DefaultCommandExecutor{}
	agentManager := ssh.NewAgentManager(executor)
	loadedKeys, err := agentManager.ListAgentKeys()
	if err != nil {
		return err
	}

	headers := []string{"NAME", "TYPE", "STATUS", "PATH"}
	var rows [][]string
	for _, privateKey := range privateKeys {
		keyName := filepath.Base(privateKey)
		keyInfo := getKeyInfo(privateKey)

		publicKey := readPublicKey(privateKey + ".pub")

		status := "Not Loaded"
		if slices.Contains(loadedKeys, publicKey) {
			status = "Loaded"
		}

		path := utils.ReplaceHomeDirWithTilde(privateKey)
		rows = append(rows, []string{keyName, keyInfo, status, path})
	}

	utils.PrintTable(os.Stdout, headers, rows)

	return nil
}

// findPrivateKeys finds all SSH private key files in the given directory
func findPrivateKeys(sshPath string) ([]string, error) {
	var privateKeys []string

	err := filepath.WalkDir(sshPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() || strings.HasSuffix(path, ".pub") {
			return nil
		}

		name := entry.Name()
		if name == "config" || name == "known_hosts" || name == "authorized_keys" {
			return nil
		}

		if isPrivateKeyFile(path) {
			privateKeys = append(privateKeys, path)
		}

		return nil
	})

	return privateKeys, err
}

// isPrivateKeyFile checks if a file is likely an SSH private key
func isPrivateKeyFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// just read first few bytes to check for private key headers
	buffer := make([]byte, 100)
	n, err := file.Read(buffer)
	if err != nil {
		return false
	}

	content := string(buffer[:n])
	return strings.Contains(content, "-----BEGIN") &&
		(strings.Contains(content, "PRIVATE KEY") || strings.Contains(content, "OPENSSH PRIVATE KEY"))
}

func getKeyInfo(keyPath string) string {
	name := filepath.Base(keyPath)

	if strings.Contains(name, "ed25519") {
		return "ED25519"
	} else if strings.Contains(name, "rsa") {
		return "RSA"
	} else if strings.Contains(name, "ecdsa") {
		return "ECDSA"
	}

	return "Unknown"
}

func readPublicKey(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		utils.PrintWarning("Failed to read public key file " + path + ": " + err.Error())
		return ""
	}

	return strings.TrimSpace(string(content))
}
