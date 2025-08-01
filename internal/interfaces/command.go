package interfaces

import (
	"os/exec"
)

type CommandExecutor interface {
	Execute(name string, args ...string) error
	ExecuteWithOutput(name string, args ...string) ([]byte, error)
}

type DefaultCommandExecutor struct{}

func (r *DefaultCommandExecutor) Execute(name string, args ...string) error {
	return exec.Command(name, args...).Run()
}

func (r *DefaultCommandExecutor) ExecuteWithOutput(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}
