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
		bin := "/engineyard/buildpacks/awesome/bin"
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

		steps.Runner.Add(detect)

		file, err = fs.Root.Create(compile)
		if err != nil {
			return err
		}
		file.Close()

		steps.Runner.Add(compile)

		return nil
	})

	s.Step(`^there is a buildpack installed that cannot build toast$`, func() error {
		bin := "/engineyard/buildpacks/onoes/bin"

		err := fs.CreateDir(bin, 0755)
		if err != nil {
			return err
		}

		file, err := fs.Root.Create(bin + "/detect")
		if err != nil {
			return err
		}
		file.Close()

		file, err = fs.Root.Create(bin + "/compile")
		if err != nil {
			return err
		}
		file.Close()

		return nil
	})

	s.Step(`^the shared config is applied to the build$`, func() error {
		actual := steps.Logger.Logged()
		expected := "Applied shared config to build"

		if !strings.Contains(actual, expected) {
			return fmt.Errorf("Expected '%s' to include '%s'", actual, expected)
		}

		return nil
	})

	s.Step(`^the proper buildpack is detected`, func() error {
		bin := "/engineyard/buildpacks/awesome/bin"
		compile := bin + "/compile"

		if !steps.Runner.Ran(compile) {
			return fmt.Errorf("Expected %s to build the app", compile)
		}

		return nil
	})

	s.Step(`^a new toast slug is generated`, func() error {
		actual := steps.Logger.Logged()
		expected := "Package saved to"

		if !strings.Contains(actual, expected) {
			return fmt.Errorf("Expected '%s' to include '%s'", actual, expected)
		}

		return nil
	})

	s.BeforeSuite(func() {
		mockable.Enable()
	})

	s.AfterSuite(func() {
		mockable.Disable()
	})

	s.BeforeScenario(func(interface{}) {
		fs.Root = afero.NewMemMapFs()
		cmd.Logger = steps.Logger
		os.NewRunner = func() ogun.Runner {
			return steps.Runner
		}
		os.NewLoggedRunner = func(context string, logger ogun.Logger) ogun.Runner {
			return steps.Runner
		}

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
