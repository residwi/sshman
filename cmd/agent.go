package cmd

import (
	"fmt"

	"github.com/residwi/sshman/internal/interfaces"
	"github.com/residwi/sshman/internal/ssh"
	"github.com/residwi/sshman/utils"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage SSH agent",
	Long:  `Add, remove, list, or clear SSH keys from agent`,
}

var agentAddCmd = &cobra.Command{
	Use:   "add [key-name]",
	Short: "Add SSH key to agent",
	Long:  `Add an SSH key to the agent for authentication`,
	Args:  cobra.ExactArgs(1),
	Example: `sshman agent add id_ed25519_work
sshman agent add id_rsa_personal`,
	RunE: addKeyToAgent,
}

var agentRemoveCmd = &cobra.Command{
	Use:   "remove [key-name]",
	Short: "Remove SSH key from agent",
	Long:  `Remove an SSH key from the agent`,
	Args:  cobra.ExactArgs(1),
	Example: `sshman agent remove id_ed25519_work
sshman agent remove id_rsa_personal`,
	RunE: removeKeyFromAgent,
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List keys loaded in agent",
	Long:  `Display all SSH keys currently loaded in agent`,
	Args:  cobra.NoArgs,
	RunE:  listAgentKeys,
}

var agentClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all keys from agent",
	Long:  `Remove all SSH keys from agent`,
	Args:  cobra.NoArgs,
	RunE:  clearAgentKeys,
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentAddCmd, agentRemoveCmd, agentListCmd, agentClearCmd)
}

func addKeyToAgent(cmd *cobra.Command, args []string) error {
	sshPath, _ := cmd.Flags().GetString("ssh-path")
	keyName := args[0]

	executor := &interfaces.DefaultCommandExecutor{}
	agentManager := ssh.NewAgentManager(executor)
	if !agentManager.IsAgentRunning() {
		return fmt.Errorf("agent is not running")
	}

	if err := agentManager.AddToAgent(sshPath, keyName); err != nil {
		return err
	}

	utils.PrintSuccess("SSH key [" + keyName + "] added to agent")
	return nil
}

func removeKeyFromAgent(cmd *cobra.Command, args []string) error {
	sshPath, _ := cmd.Flags().GetString("ssh-path")
	keyName := args[0]

	executor := &interfaces.DefaultCommandExecutor{}
	agentManager := ssh.NewAgentManager(executor)

	if !agentManager.IsAgentRunning() {
		return fmt.Errorf("agent is not running")
	}

	if err := agentManager.RemoveFromAgent(sshPath, keyName); err != nil {
		return err
	}

	utils.PrintSuccess("SSH key [" + keyName + "] removed from agent")
	return nil
}

func listAgentKeys(cmd *cobra.Command, args []string) error {
	executor := &interfaces.DefaultCommandExecutor{}
	agentManager := ssh.NewAgentManager(executor)

	keys, err := agentManager.ListAgentKeys()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		utils.PrintSuccess("No SSH keys loaded in agent")
		return nil
	}

	for _, key := range keys {
		fmt.Println(key)
	}

	return nil
}

func clearAgentKeys(cmd *cobra.Command, args []string) error {
	sshPath, _ := cmd.Flags().GetString("ssh-path")

	executor := &interfaces.DefaultCommandExecutor{}
	agentManager := ssh.NewAgentManager(executor)

	if !agentManager.IsAgentRunning() {
		return fmt.Errorf("agent is not running")
	}

	if err := agentManager.ClearAgent(sshPath); err != nil {
		return err
	}

	return nil
}
