package main

import (
	"os"

	"github.com/ess/conan/cmd/conan/cmd"
)

func main() {
	if cmd.Execute() != nil {
		os.Exit(1)
	}
}
