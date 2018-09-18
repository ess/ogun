package os

import (
	"os/exec"
)

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

func (runner *Runner) Execute(command string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", command)

	output, err := cmd.CombinedOutput()
	return output, err
}
