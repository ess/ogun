package fs

import (
	"github.com/ess/ogun"
)

type ApplicationService struct {
	logger ogun.Logger
}

func NewApplicationService(logger ogun.Logger) ApplicationService {
	return ApplicationService{logger: logger}
}

func (service ApplicationService) Get(name string) (ogun.Application, error) {
	return ogun.Application{Name: name}, nil
}
