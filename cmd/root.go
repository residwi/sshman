package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/residwi/sshman/utils"
	"github.com/spf13/cobra"
)

var rootCmdFlags struct {
	sshPath string
}

var rootCmd = &cobra.Command{
	Use:   "sshman",
	Short: "manage SSH keys",
	Long:  `sshman is a CLI tool to manage SSH keys and config easily.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		sshPath, _ := cmd.Flags().GetString("ssh-path")
		if utils.IsDirectoryNotExist(sshPath) {
			return fmt.Errorf("SSH path does not exist: %s", sshPath)
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.PrintError(err.Error())
		os.Exit(1)
	}
}

func init() {
	homeDir, _ := os.UserHomeDir()
	sshDefaultPath := filepath.Join(homeDir, ".ssh")

	rootCmd.PersistentFlags().StringVar(&rootCmdFlags.sshPath, "ssh-path", sshDefaultPath, "Path to SSH directory")
}

func isSupportedKeyType(keyType string) error {
	supportedTypes := []string{"ed25519", "rsa"}
	if !slices.Contains(supportedTypes, keyType) {
		return fmt.Errorf("unsupported SSH key type: %s. Supported types are: %v", keyType, supportedTypes)
	}
	return nil
}
