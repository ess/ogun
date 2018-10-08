package steps

import (
	"fmt"
	"strings"

	//"github.com/ess/jamaica"
	"github.com/ess/kennel"
	"github.com/ess/mockable"
	"github.com/spf13/afero"

	"github.com/ess/ogun/cmd/ogun/cmd"
	"github.com/ess/ogun/cmd/ogun/workflows"
	"github.com/ess/ogun/mock"
	"github.com/ess/ogun/pkg/ogun"
	"github.com/ess/ogun/pkg/ogun/fs"
	"github.com/ess/ogun/pkg/ogun/os"
)

type BuildingAnApp struct {
	Logger      *mock.Logger
	Runner      *mock.Runner
	buildNumber string
}

func (steps *BuildingAnApp) logContains(expected string) error {
	actual := steps.Logger.Logged()

	if !strings.Contains(actual, expected) {
		return fmt.Errorf("expected '%s' to include '%s'", actual, expected)
	}

	return nil
}

func (steps *BuildingAnApp) stubBuildpack(base string, detectable bool, compilable bool) error {
	bin := base + "/bin"
	detect := bin + "/detect"
	compile := bin + "/compile"

	err := fs.CreateDir(bin, 0755)
	if err != nil {
		return err
	}

	file, err := fs.Root.Create(detect)
	if err != nil {
		return err
	}
	file.Close()

	if detectable {
		err = fs.Root.Chmod(detect, 0755)
		if err != nil {
			return err
		}

		steps.Runner.Add(detect)
	}

	file, err = fs.Root.Create(compile)
	if err != nil {
		return err
	}
	file.Close()

	if compilable {
		err = fs.Root.Chmod(compile, 0755)
		if err != nil {
			return err
		}

		steps.Runner.Add(compile)
	}

	return nil

}

func (steps *BuildingAnApp) StepUp(s kennel.Suite) {
	s.Step(`^there is an app named toast$`, func() error {
		return fs.CreateDir("/data/toast", 0755)
	})

	s.Step(`^there is a shared configuration for the toast app$`, func() error {
		return fs.CreateDir("/data/toast/shared/config", 0755)
	})

	s.Step(`^I have a cached copy of the toast app$`, func() error {
		return fs.CreateDir("/data/toast/shared/cached-copy", 0755)
	})

	s.Step(`^there is a buildpack installed that can build toast$`, func() error {
		return steps.stubBuildpack("/engineyard/buildpacks/awesome", true, true)
	})

	s.Step(`^there is a buildpack installed that cannot build toast$`, func() error {
		return steps.stubBuildpack("/engineyard/buildpacks/onoes", false, false)
	})

	s.Step(`^the shared config is applied to the build$`, func() error {
		return steps.logContains("Applied shared config to build")
	})

	s.Step(`^the proper buildpack is detected$`, func() error {
		bin := "/engineyard/buildpacks/awesome/bin"
		compile := bin + "/compile"

		if !steps.Runner.Ran(compile) {
			return fmt.Errorf("Expected %s to build the app", compile)
		}

		return nil
	})

	s.Step(`^a new toast slug is generated$`, func() error {
		return steps.logContains(
			"Package saved to /data/toast/releases/" + steps.buildNumber + ".tgz",
		)
	})

	s.Step(`^a new toast slug named 0987654321 is generated$`, func() error {
		return steps.logContains(
			"Package saved to /data/toast/releases/0987654321.tgz",
		)
	})

	s.Step(`^the toast app has a custom buildpack$`, func() error {
		return steps.stubBuildpack(
			"/data/toast/shared/cached-copy/.ogun/buildpack",
			true,
			true,
		)
	})

	s.Step(`^the custom buildpack is used to build the release$`, func() error {
		compile := "/data/toast/shared/cached-copy/.ogun/buildpack/bin/compile"

		if !steps.Runner.Ran(compile) {
			return fmt.Errorf("Expected %s to build the app", compile)
		}

		return nil
	})

	s.Step(`^there are no buildpacks that can build the app$`, func() error {
		fs.Root.RemoveAll("/engineyard/buildpacks")
		fs.Root.RemoveAll("/data/toast/shared/cached-copy/.ogun")
		return nil
	})

	s.Step(`^I see an error regarding the lack of a viable buildpack$`, func() error {
		return steps.logContains("Detected no buildpacks that can build toast")
	})

	s.Step(`^there is an issue building the application$`, func() error {
		steps.Runner.Remove("/engineyard/buildpacks/awesome/bin/compile")

		//fmt.Println("build ReleaseName:", cmd.ReleaseName)
		//// For some reason, if I don't explicitly reset the logger here, it holds
		//// onto a previous run's logs
		//steps.Logger.Reset()

		return nil
	})

	s.Step(`^I see an error regarding the build failure$`, func() error {
		return steps.logContains(
			"Compiling release " + steps.buildNumber + " failed",
		)
	})

	s.BeforeSuite(func() {
		mockable.Enable()
	})

	s.AfterSuite(func() {
		mockable.Disable()
	})

	s.BeforeScenario(func(interface{}) {
		steps.Logger.Reset()
		steps.Runner.Reset()

		fs.Root = afero.NewMemMapFs()
		cmd.Logger = steps.Logger
		os.NewRunner = func() ogun.Runner {
			return steps.Runner
		}
		os.NewLoggedRunner = func(context string, logger ogun.Logger) ogun.Runner {
			return steps.Runner
		}

		cmd.ReleaseName = ""

		workflows.GenerateBuildNumber = func() string {
			return steps.buildNumber
		}
	})

}

func init() {
	kennel.Register(
		&BuildingAnApp{
			Logger:      mock.NewLogger(),
			Runner:      mock.NewRunner(),
			buildNumber: "1234567890",
		},
	)
}
