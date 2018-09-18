package bash

import (
	"os/exec"
	"syscall"
)

var Run = func(commnad string) int {
	cmd := exec.Command("bash", "-c", command)
	var waitStatus syscall.WaitStatus

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			return waitStatus.ExitStatus()
		}
	}

	return 0
}
