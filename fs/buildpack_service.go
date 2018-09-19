package fs

import (
	"fmt"

	"github.com/ess/ogun"
	"github.com/ess/ogun/os"
)

type BuildpackService struct {
	logger ogun.Logger
}

func NewBuildpackService(logger ogun.Logger) BuildpackService {
	return BuildpackService{logger: logger}
}

func (service BuildpackService) Get(name string) (ogun.Buildpack, error) {
	context := "buildpack-get"
	path := buildpackRoot() + name
	detect := path + "/bin/detect"
	compile := path + "/bin/compile"

	if !FileExists(path) || !FileExists(detect) || !FileExists(compile) {
		message := name + " is not a valid buildpack"

		service.logger.Error(context, message)

		return ogun.Buildpack{}, fmt.Errorf(message)
	}

	return ogun.Buildpack{Name: name}, nil
}

func (service BuildpackService) Detect(application ogun.Application) (ogun.Buildpack, error) {
	context := "buildpack-detect"
	detected := make([]ogun.Buildpack, 0)

	for _, pack := range service.all() {
		err := service.detect(application, pack)
		if err == nil {
			service.logger.Info(
				context,
				pack.Name+" supports building this application",
			)

			detected = append(detected, pack)
		}
	}

	if len(detected) < 1 {
		service.logger.Error(
			context,
			"Detected no buildpacks that can build "+application.Name,
		)

		return ogun.Buildpack{}, fmt.Errorf("No buildpacks support this application")
	}

	if len(detected) > 1 {
		service.logger.Error(
			context,
			"Detected multiple buildpacks that can build "+application.Name,
		)

		return ogun.Buildpack{}, fmt.Errorf("Multiple buildpacks support this application")
	}

	return detected[0], nil
}

func (service BuildpackService) all() []ogun.Buildpack {
	buildpacks := make([]ogun.Buildpack, 0)

	if candidates, err := ReadDir(buildpackRoot()); err == nil {

		for _, info := range candidates {
			name := Basename(info.Name())
			buildpacks = append(buildpacks, ogun.Buildpack{Name: name})
		}
	}

	return buildpacks
}

func (service BuildpackService) detect(app ogun.Application, pack ogun.Buildpack) error {
	detectPath := buildpackPath(pack) + "/bin/detect"
	cacheRoot := applicationPath(app) + "/shared/cached_copy"

	_, err := os.NewRunner().Execute(detectPath + " " + cacheRoot)

	return err
}
