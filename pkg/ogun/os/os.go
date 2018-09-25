package os

import (
	"github.com/ess/ogun/pkg/ogun"
)

func varsToStrings(vars []ogun.Variable) []string {
	output := make([]string, 0)

	for _, variable := range vars {
		output = append(output, string(variable))
	}

	return output
}
