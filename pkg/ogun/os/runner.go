package os

import (
	"os/exec"

	"github.com/ess/ogun/pkg/ogun"
)

type Runner struct{}

var NewRunner = func() ogun.Runner {
	return &Runner{}
}

func (runner *Runner) Execute(command string, vars []ogun.Variable) ([]byte, error) {
	cmd := exec.Command("bash", "-c", command)

	if len(vars) > 0 {
		cmd.Env = varsToStrings(vars)
	}

	output, err := cmd.CombinedOutput()
	return output, err
}
