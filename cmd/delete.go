package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/residwi/sshman/internal/interfaces"
	"github.com/residwi/sshman/internal/ssh"
	"github.com/residwi/sshman/utils"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [key-name]",
	Short: "Delete SSH key and remove from agent",
	Long:  `Delete an SSH key pair and remove it from agent`,
	Args:  cobra.ExactArgs(1),
	Example: `sshman delete id_ed25519_work
sshman delete id_rsa_personal`,
	RunE: deleteSSHKey,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteSSHKey(cmd *cobra.Command, args []string) error {
	sshPath, _ := cmd.Flags().GetString("ssh-path")
	keyName := args[0]

	executor := &interfaces.DefaultCommandExecutor{}
	agentManager := ssh.NewAgentManager(executor)
	if agentManager.IsAgentRunning() {
		if err := agentManager.RemoveFromAgent(sshPath, keyName); err == nil {
			utils.PrintSuccess("SSH key [" + keyName + "] removed from agent")
		}
	}

	privateKeyPath := filepath.Join(sshPath, keyName)
	if err := os.Remove(privateKeyPath); err != nil {
		return fmt.Errorf("failed to delete private key: %w", err)
	}
	utils.PrintSuccess("SSH private key [" + keyName + "] deleted successfully")

	publicKeyPath := privateKeyPath + ".pub"
	if err := os.Remove(publicKeyPath); err != nil {
		return fmt.Errorf("failed to delete public key: %w", err)
	}
	utils.PrintSuccess("SSH public key [" + keyName + "] deleted successfully")
	utils.PrintWarning("Note: You may need to manually remove the key from your SSH config file")

	return nil
}
