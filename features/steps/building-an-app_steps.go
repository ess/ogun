package steps

import (
	//"fmt"
	//"strings"

	//"github.com/ess/jamaica"
	"github.com/ess/kennel"
	"github.com/ess/mockable"
	"github.com/spf13/afero"

	"github.com/ess/ogun/cmd/ogun/cmd"
	"github.com/ess/ogun/mock"
	"github.com/ess/ogun/pkg/ogun/fs"
)

type BuildingAnApp struct{}

func (steps *BuildingAnApp) StepUp(s kennel.Suite) {
	s.Step(`^there is an app named toast$`, func() error {
		return fs.CreateDir("/data/toast", 0755)
	})

	s.Step(`^there is a shared configuration for the toast app$`, func() error {
		return fs.CreateDir("/data/toast/shared/config", 0755)
	})

	s.Step(`^I have a cached copy of the toast app$`, func() error {
		return fs.CreateDir("/data/toast/shared/cached_copy", 0755)
	})

	s.Step(`^there is a buildpack installed that can build toast$`, func() error {
		bin := "/engineyard/buildpacks/awesome/bin"

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

	//s.Step(`^the shared config is applied to the build$`, func() error {
	//if !strings.Contains(jamaica.LastCommandStdout())
	//})

	s.BeforeSuite(func() {
		mockable.Enable()
	})

	s.AfterSuite(func() {
		mockable.Disable()
	})

	s.BeforeScenario(func(interface{}) {
		fs.Root = afero.NewMemMapFs()
		cmd.Logger = mock.NewLogger()
	})

}

func init() {
	kennel.Register(new(BuildingAnApp))
}
