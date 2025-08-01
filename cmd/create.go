package cmd

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/residwi/sshman/internal/interfaces"
	"github.com/residwi/sshman/internal/provider"
	"github.com/residwi/sshman/internal/ssh"
	"github.com/residwi/sshman/utils"
	"github.com/spf13/cobra"
)

const (
	defaultSSHKeyAlgorithm = "ed25519"
)

var createCmdFlags struct {
	typeKey  string
	email    string
	purpose  string
	user     string
	hostname string
}

var createCmd = &cobra.Command{
	Use:       "create [github|gitlab|bitbucket|generic]",
	Short:     "Create a new SSH key",
	ValidArgs: provider.GetSupportedProviders(),
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Example:   `sshman create github --email residwi@mail.com -t ed25519 --purpose work`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := isSupportedKeyType(createCmdFlags.typeKey); err != nil {
			return err
		}

		if args[0] == "generic" && (createCmdFlags.user == "" || createCmdFlags.hostname == "") {
			return fmt.Errorf("for 'generic' provider, both --user and --hostname flags are required")
		}

		if slices.Contains(provider.GetSupportedProviders(), args[0]) {
			if createCmdFlags.user != "" || createCmdFlags.hostname != "" {
				utils.PrintWarning("for provider " + args[0] + ", --user and --hostname flags are ignored. Using default values from provider config")
			}

			if createCmdFlags.email == "" {
				return fmt.Errorf("email is required for provider %s. Use --email flag", args[0])
			}
		}

		return nil
	},
	RunE: generateSSH,
}

func init() {
	typeKeys := strings.Join([]string{defaultSSHKeyAlgorithm, "rsa"}, ", ")
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createCmdFlags.typeKey, "type", "t", defaultSSHKeyAlgorithm, "Type of the SSH key ("+typeKeys+")")
	createCmd.Flags().StringVarP(&createCmdFlags.email, "email", "", "", "Email for the public key (required)")
	createCmd.Flags().StringVarP(&createCmdFlags.purpose, "purpose", "", "", "Purpose of the SSH key (work, personal, etc.)")
	createCmd.Flags().StringVarP(&createCmdFlags.user, "user", "", "", "Username for the SSH key (generic only)")
	createCmd.Flags().StringVarP(&createCmdFlags.hostname, "hostname", "H", "", "Hostname for the SSH key (generic only)")
}

func generateSSH(cmd *cobra.Command, args []string) error {
	executor := &interfaces.DefaultCommandExecutor{}
	keyGen := ssh.NewKeyGenerator(executor)

	providerConfig, exists := provider.GetProviderConfig(args[0])
	if exists {
		createCmdFlags.user = providerConfig.User
		createCmdFlags.hostname = providerConfig.Hostname
	}

	keyConfig := ssh.KeyConfig{
		Type:     createCmdFlags.typeKey,
		Email:    createCmdFlags.email,
		Purpose:  createCmdFlags.purpose,
		Provider: args[0],
		SSHPath:  rootCmdFlags.sshPath,
	}

	keyName, err := keyGen.GenerateKey(keyConfig)
	if err != nil {
		return err
	}

	utils.PrintSuccess("SSH key [" + keyName + "] created!")

	hostAlias := getHostAlias(args[0], createCmdFlags.hostname, createCmdFlags.purpose)
	configEntry := ssh.ConfigEntry{
		Host:         hostAlias,
		User:         createCmdFlags.user,
		Hostname:     createCmdFlags.hostname,
		IdentityFile: filepath.Join(rootCmdFlags.sshPath, keyName),
	}

	if err := ssh.AddToConfig(rootCmdFlags.sshPath, configEntry); err != nil {
		return err
	}

	utils.PrintSuccess("SSH config added for host [" + hostAlias + "] with user [" + createCmdFlags.user + "]")

	executor = &interfaces.DefaultCommandExecutor{}
	agentManager := ssh.NewAgentManager(executor)
	if agentManager.IsAgentRunning() {
		if err := agentManager.AddToAgent(rootCmdFlags.sshPath, keyName); err != nil {
			utils.PrintWarning("Warning: Failed to add key to ssh-agent: " + err.Error())
			utils.PrintWarning("You can manually add the key using: ssh-add " + filepath.Join(rootCmdFlags.sshPath, keyName))
		} else {
			utils.PrintSuccess("SSH key automatically added to ssh-agent")
		}
	}

	return nil
}

func getHostAlias(provider, hostname, purpose string) string {
	if purpose != "" {
		return provider + "-" + purpose
	}
	return hostname
}
