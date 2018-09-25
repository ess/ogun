package fs

import (
	"fmt"

	"github.com/ess/ogun/pkg/ogun"
	"github.com/ess/ogun/pkg/ogun/os"
)

type BuildpackService struct {
	runner ogun.Runner
	logger ogun.Logger
}

func NewBuildpackService(logger ogun.Logger) BuildpackService {
	return BuildpackService{logger: logger, runner: os.NewRunner()}
}

func (service BuildpackService) Detect(application ogun.Application) (ogun.Buildpack, error) {
	context := "buildpack-detect"

	custom, customErr := service.custom(application)
	if customErr == nil {
		return custom, nil
	}

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

func (service BuildpackService) validate(pack ogun.Buildpack) error {
	base := pack.Location

	detect := base + "/bin/detect"
	compile := base + "/bin/compile"

	for _, path := range []string{base, detect, compile} {
		if !FileExists(path) {
			return fmt.Errorf("%s does not exist", path)
		}
	}

	for _, path := range []string{detect, compile} {
		if !Executable(path) {
			return fmt.Errorf("%s is not executable", path)
		}
	}

	return nil
}

func (service BuildpackService) custom(application ogun.Application) (ogun.Buildpack, error) {
	context := "detect-custom"

	service.logger.Info(context, "Checking for a custom buildpack ...")
	location := applicationPath(application) + "/shared/cached_copy/.ogun/buildpack"

	pack := ogun.Buildpack{Name: "custom", Location: location}

	err := service.validate(pack)
	if err != nil {
		service.logger.Info(context, "Custom buildpack found!")
		return pack, err
	}

	service.logger.Info(context, "Didn't find a valid custom buildpack")
	return pack, nil

}

func (service BuildpackService) all() []ogun.Buildpack {
	buildpacks := make([]ogun.Buildpack, 0)
	root := buildpackRoot()

	if candidates, err := ReadDir(root); err == nil {

		for _, info := range candidates {
			name := Basename(info.Name())
			buildpacks = append(buildpacks, ogun.Buildpack{Name: name, Location: root + "/" + name})
		}
	}

	return buildpacks
}

func (service BuildpackService) detect(app ogun.Application, pack ogun.Buildpack) error {
	detectPath := buildpackPath(pack) + "/bin/detect"
	cacheRoot := applicationPath(app) + "/shared/cached_copy"

	_, err := service.runner.Execute(detectPath+" "+cacheRoot, nil)

	return err
}
