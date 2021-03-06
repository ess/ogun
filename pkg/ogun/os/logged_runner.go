package os

import (
	"bytes"
	"io"
	"os/exec"
	"strings"

	"github.com/ess/ogun/pkg/ogun"
)

type LoggedRunner struct {
	context string
	logger  ogun.Logger
}

var NewLoggedRunner = func(context string, logger ogun.Logger) ogun.Runner {
	return &LoggedRunner{context: context, logger: logger}
}

func (runner *LoggedRunner) Execute(command string, vars []ogun.Variable) ([]byte, error) {
	cmd := exec.Command("bash", "-c", command)

	if len(vars) > 0 {
		cmd.Env = varsToStrings(vars)
	}

	output := make([]byte, 0)
	buf := bytes.NewBuffer(output)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	stdout := newPassThrough(runner.logger, runner.context, "info", buf)
	stderr := newPassThrough(runner.logger, runner.context, "error", buf)

	cmd.Start()

	go func() {
		io.Copy(stdout, stdoutIn)
	}()

	go func() {
		io.Copy(stderr, stderrIn)
	}()

	err := cmd.Wait()

	return buf.Bytes(), err
}

type passThrough struct {
	log     ogun.Logger
	context string
	level   string
	output  *bytes.Buffer
}

func newPassThrough(log ogun.Logger, context string, level string, output *bytes.Buffer) *passThrough {
	return &passThrough{
		log:     log,
		context: context,
		level:   level,
		output:  output,
	}
}

func (p *passThrough) Write(d []byte) (int, error) {
	p.output.Write(d)

	//line := strings.TrimSpace(string(d))
	lines := strings.Split(string(d), "\n")

	for _, line := range lines {
		if len(line) > 0 {
			switch p.level {
			case "error":
				p.log.Error(p.context, line)
			default:
				p.log.Info(p.context, line)
			}
		}
	}

	return len(d), nil
}
