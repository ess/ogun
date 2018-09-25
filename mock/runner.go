package mock

import (
	"fmt"
	"strings"

	"github.com/ess/ogun/pkg/ogun"
)

type Runner struct {
	goodPrefixes []string
}

func NewRunner() *Runner {
	return &Runner{}
}

func (runner *Runner) Execute(command string, vars []ogun.Variable) ([]byte, error) {
	if strings.HasSuffix(command, " env") {
		env := make([]string, 0)
		for _, variable := range vars {
			env = append(env, string(variable))
		}

		return []byte(strings.Join(env, "\n")), nil
	}

	return []byte(fmt.Sprintf("RUNNER: %s", command)), runner.generateError(command)
}

func (runner *Runner) Add(prefix string) {
	runner.goodPrefixes = append(runner.goodPrefixes, prefix)
}

func (runner *Runner) Remove(command string) {
	index := runner.find(command)

	if index >= 0 {
		runner.goodPrefixes = append(runner.goodPrefixes[:index], runner.goodPrefixes[index+1:]...)
	}
}

func (runner *Runner) Reset() {
	runner.goodPrefixes = make([]string, 0)
}

func (runner *Runner) generateError(command string) error {
	if runner.found(command) {
		return nil
	}

	return fmt.Errorf("command failed")
}

func (runner *Runner) find(command string) int {
	for index, candidate := range runner.goodPrefixes {
		if strings.HasPrefix(command, candidate) {
			return index
		}
	}

	return -1
}

func (runner *Runner) found(command string) bool {
	if runner.find(command) >= 0 {
		return true
	}

	return false
}
