package fs

import (
	"fmt"

	"github.com/ess/ogun/pkg/ogun"
)

type ApplicationService struct {
	logger ogun.Logger
}

func NewApplicationService(logger ogun.Logger) ApplicationService {
	return ApplicationService{logger: logger}
}

func (service ApplicationService) Get(name string) (ogun.Application, error) {
	app := ogun.Application{Name: name}

	if !DirectoryExists(applicationPath(app)) {
		return app, fmt.Errorf("no application named %s found", name)
	}

	return app, nil
}
