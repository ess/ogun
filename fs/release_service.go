package fs

import (
	"fmt"

	"github.com/ess/ogun"
	exec "github.com/ess/ogun/os"
)

type ReleaseService struct {
	Logger ogun.Logger
}

func NewReleaseService(logger ogun.Logger) ReleaseService {
	return ReleaseService{Logger: logger}
}

func (service ReleaseService) Create(name string, app ogun.Application) (ogun.Release, error) {
	rel := ogun.Release{Name: name, Application: app}
	path, err := service.createReleasePath(app, name)
	if err != nil {
		return rel, fmt.Errorf("could not create " + path)
	}

	err = service.copySource(app, path)

	return rel, err
}

func (service ReleaseService) applyConfig(release ogun.Release) error {
	context := "apply-config"
	app := release.Application

	source := configPath(app)
	destination := applicationPath(app) + "/builds/" + release.Name + "/config"

	err := DirCopy(source, destination)
	if err != nil {
		service.Logger.Error(context, "Could not apply shared config to build")
	} else {
		service.Logger.Info(context, "Applied shared config to build")
	}

	return err
}

func (service ReleaseService) Build(release ogun.Release, pack ogun.Buildpack) error {
	compile := buildpackPath(pack) + "/bin/compile"
	buildPath := applicationPath(release.Application) + "/builds/" + release.Name

	// applyConfig or bail
	err := service.applyConfig(release)
	if err != nil {
		return err
	}

	runner := exec.NewLoggedRunner(pack.Name+"/compile", service.Logger)
	_, err = runner.Execute(compile + " " + buildPath + " " + cachePath(release.Application))

	if err != nil {
		service.Logger.Error("build", "Compiling release "+release.Name+" failed")
	}

	return err
}

func (service ReleaseService) Clean(release ogun.Release) error {
	return fmt.Errorf("not implemented")
}

func (service ReleaseService) Package(release ogun.Release) error {
	context := "release/package"

	buildPath := applicationPath(release.Application) + "/builds/" + release.Name
	slugPath := applicationPath(release.Application) + "/releases/"
	slugFile := slugPath + release.Name + ".tgz"

	service.Logger.Info(context, "Packaging "+release.Name)

	err := CreateDir(slugPath, 0755)
	if err != nil {
		service.Logger.Error(context, "Release storage directory does not exist")
		return err
	}

	err = Tar(buildPath, slugFile)
	if err != nil {
		service.Logger.Error(context, "Could not create package for "+release.Name)
		return err
	}

	service.Logger.Info(context, "Package saved to "+slugFile)
	return nil
}

func (service ReleaseService) createReleasePath(app ogun.Application, name string) (string, error) {
	path := applicationPath(app) + "/builds/" + name

	err := CreateDir(path, 0700)

	return path, err
}

func (service ReleaseService) copySource(app ogun.Application, path string) error {

	cacheRoot := applicationPath(app) + "/shared/cached_copy"

	return DirCopy(cacheRoot, path)
}

func (service ReleaseService) Delete(release ogun.Release) error {
	return fmt.Errorf("not implemented")
}
