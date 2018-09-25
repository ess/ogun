package main

import (
	"os"

	"github.com/ess/ogun/cmd/ogun/cmd"
)

func main() {
	if cmd.Execute() != nil {
		os.Exit(1)
	}
}
