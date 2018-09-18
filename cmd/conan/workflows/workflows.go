package workflows

import (
	"fmt"
)

type Workflow interface {
	perform() error
}

func Perform(wf Workflow) error {
	return wf.perform()
}

func fatality() error {
	return fmt.Errorf("Cannot continue")
}
