package fs

import (
	"github.com/ess/conan"
)

type ApplicationService struct {
	logger conan.Logger
}

func NewApplicationService(logger conan.Logger) ApplicationService {
	return ApplicationService{logger: logger}
}

func (service ApplicationService) Get(name string) (conan.Application, error) {
	return conan.Application{Name: name}, nil
}
