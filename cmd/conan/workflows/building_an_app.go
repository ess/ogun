package workflows

import (
	"fmt"
	"time"

	"github.com/ess/conan"
)

type BuildingAnApp struct {
	BuildpackName   string
	ApplicationName string
	ReleaseName     string
	Apps            conan.ApplicationService
	Packs           conan.BuildpackService
	Releases        conan.ReleaseService
	Logger          conan.Logger
}

func (workflow *BuildingAnApp) perform() error {
	context := "main"
	workflow.Logger.Info(context, "Starting ...")

	app, err := workflow.loadApplication()
	if err != nil {
		return fatality()
	}

	pack, err := workflow.loadBuildpack(app)
	if err != nil {
		return fatality()
	}

	release, err := workflow.createRelease(app)
	if err != nil {
		return fatality()
	}

	workflow.Logger.Info(
		context,
		fmt.Sprintf(
			"Building release %s for %s (%s)",
			release.Name,
			app.Name,
			pack.Name,
		),
	)

	err = workflow.Releases.Build(release, pack)
	if err != nil {
		return fatality()
	}

	err = workflow.Releases.Package(release)
	if err != nil {
		return fatality()
	}

	return nil
}

func (workflow *BuildingAnApp) loadApplication() (conan.Application, error) {
	app, err := workflow.Apps.Get(workflow.ApplicationName)
	if err != nil {
		workflow.Logger.Error(
			"load-app",
			"Could not find application "+workflow.ApplicationName,
		)
	}

	return app, err
}

func (workflow *BuildingAnApp) loadBuildpack(app conan.Application) (conan.Buildpack, error) {
	var pack conan.Buildpack
	var err error

	if workflow.BuildpackName == "detect" {
		pack, err = workflow.Packs.Detect(app)
	} else {
		pack, err = workflow.Packs.Get(workflow.BuildpackName)
	}

	if err != nil {
		workflow.Logger.Error(
			"load-buildpack",
			"Could not detect buildpack for "+app.Name,
		)
	}

	return pack, err
}

func (workflow *BuildingAnApp) createRelease(app conan.Application) (conan.Release, error) {

	release, err := workflow.Releases.Create(workflow.ReleaseName, app)
	if err != nil {
		fmt.Println("create error:", err.Error())
		workflow.Logger.Error(
			"create-release",
			"Could not create release "+release.Name+" for "+app.Name,
		)

		return release, err
	}

	return release, nil
}

func GenerateBuildNumber() string {
	now := time.Now()

	return fmt.Sprintf(
		"%d%02d%02d%02d%02d%02d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)
}
