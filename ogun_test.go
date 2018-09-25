package main

import (
	"os"
	"testing"
	"time"

	"github.com/ess/ogun/features/steps"

	"github.com/DATA-DOG/godog"
	"github.com/ess/jamaica"
	"github.com/ess/kennel"
	"github.com/ess/mockable"
	"github.com/ess/testscope"

	"github.com/ess/ogun/cmd/ogun/cmd"
)

var commandOutput string
var lastCommandRanErr error

func TestMain(m *testing.M) {
	if testscope.Integration() {
		status := godog.RunWithOptions(
			"godog",
			func(s *godog.Suite) {
				mockable.Enable()
				jamaica.SetRootCmd(cmd.RootCmd)
				jamaica.StepUp(s)
				steps.Register()
				kennel.StepUp(s)
			},

			godog.Options{
				Format:    "pretty",
				Paths:     []string{"features"},
				Randomize: time.Now().UTC().UnixNano(),
			},
		)

		if st := m.Run(); st > status {
			status = st
		}

		os.Exit(status)
	}
}

func TestTrue(t *testing.T) {
}
