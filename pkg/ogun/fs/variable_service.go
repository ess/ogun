package fs

import (
	"os"
	"sort"
	"strings"

	"github.com/ess/ogun/pkg/ogun"
	myos "github.com/ess/ogun/pkg/ogun/os"
)

type VariableService struct {
	logger ogun.Logger
}

func NewVariableService(logger ogun.Logger) VariableService {
	return VariableService{logger: logger}
}

func (service VariableService) All(app ogun.Application) []ogun.Variable {
	envs := service.stringsToVars(os.Environ())

	runner := myos.NewRunner()

	for _, path := range service.envfiles(app) {
		aggregate, err := runner.Execute("source "+path+" ; env", envs)
		if err != nil {
			continue
		}

		//envs = []ogun.Variable(strings.Split(string(aggregate), "\n"))
		envs = service.stringsToVars(strings.Split(string(aggregate), "\n"))
	}

	return envs
}

func (service VariableService) envfiles(app ogun.Application) []string {
	basedir := configPath(app)
	standard := []string{
		basedir + "/env",
		basedir + "/env.cloud",
		basedir + "/env.custom",
	}

	uniques := make(map[string]bool)

	for _, file := range standard {
		uniques[file] = true
	}

	matches, _ := ReadDir(basedir)
	for _, match := range matches {
		name := match.Name()

		if strings.HasPrefix(name, "env.") {
			uniques[basedir+"/"+name] = true
		}
	}

	for _, file := range standard {
		delete(uniques, file)
	}

	nonstandard := make([]string, 0)

	for file := range uniques {
		nonstandard = append(nonstandard, file)
	}

	sort.Strings(nonstandard)

	candidates := append(standard, nonstandard...)

	envfiles := make([]string, 0)

	for _, candidate := range candidates {
		if FileExists(candidate) {
			envfiles = append(envfiles, candidate)
		}
	}

	return envfiles
}

func (service VariableService) stringsToVars(vars []string) []ogun.Variable {
	output := make([]ogun.Variable, 0)

	for _, variable := range vars {
		output = append(output, ogun.Variable(variable))
	}

	return output

}
