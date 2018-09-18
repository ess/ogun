package fs

import (
	"fmt"

	"github.com/ess/conan"
	"github.com/ess/conan/os"
)

type BuildpackService struct {
	logger conan.Logger
}

func NewBuildpackService(logger conan.Logger) BuildpackService {
	return BuildpackService{logger: logger}
}

func (service BuildpackService) Get(name string) (conan.Buildpack, error) {
	context := "buildpack-get"
	path := buildpackRoot() + name
	detect := path + "/bin/detect"
	compile := path + "/bin/compile"

	if !FileExists(path) || !FileExists(detect) || !FileExists(compile) {
		message := name + " is not a valid buildpack"

		service.logger.Error(context, message)

		return conan.Buildpack{}, fmt.Errorf(message)
	}

	return conan.Buildpack{Name: name}, nil
}

func (service BuildpackService) Detect(application conan.Application) (conan.Buildpack, error) {
	context := "buildpack-detect"
	detected := make([]conan.Buildpack, 0)

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

		return conan.Buildpack{}, fmt.Errorf("No buildpacks support this application")
	}

	if len(detected) > 1 {
		service.logger.Error(
			context,
			"Detected multiple buildpacks that can build "+application.Name,
		)

		return conan.Buildpack{}, fmt.Errorf("Multiple buildpacks support this application")
	}

	return detected[0], nil
}

func (service BuildpackService) all() []conan.Buildpack {
	buildpacks := make([]conan.Buildpack, 0)

	if candidates, err := ReadDir(buildpackRoot()); err == nil {

		for _, info := range candidates {
			name := Basename(info.Name())
			buildpacks = append(buildpacks, conan.Buildpack{Name: name})
		}
	}

	return buildpacks
}

func (service BuildpackService) detect(app conan.Application, pack conan.Buildpack) error {
	detectPath := buildpackPath(pack) + "/bin/detect"
	cacheRoot := applicationPath(app) + "/shared/cached_copy"

	_, err := os.NewRunner().Execute(detectPath + " " + cacheRoot)

	return err
}
