package fs

import (
	"github.com/brianvoe/gofakeit"
	"github.com/spf13/afero"

	"github.com/ess/ogun/mock"
	"github.com/ess/ogun/pkg/ogun"
)

var logger = mock.NewLogger()
var appName string
var releaseName = "20180921143528"
var app ogun.Application
var release ogun.Release
var runner = mock.NewRunner()

func setupFs() {
	Root = afero.NewMemMapFs()
}

func stubApplication() {
	appName = gofakeit.Generate("????????????")
	app = ogun.Application{Name: appName}

	CreateDir(applicationPath(app), 0755)
	CreateDir(applicationPath(app)+"/shared/cached-copy", 0755)

	ohaipath := applicationPath(app) + "/shared/cached-copy/ohai"
	ohai, _ := Root.Create(ohaipath)
	ohai.Close()

	if !FileExists(ohaipath) {
		panic("what the fuck man")
	}
}

func stubPack(pack ogun.Buildpack) error {
	bin := pack.Location + "/bin"

	CreateDir(bin, 0755)

	detect := bin + "/detect"
	//fmt.Println("generating", detect)
	file, err := Root.Create(detect)
	if err != nil {
		return err
	}
	file.Close()
	Root.Chmod(detect, 0755)
	//fmt.Println(detect, "exists?", FileExists(detect))
	//fmt.Println(detect, "executable?", Executable(detect))

	compile := bin + "/compile"
	//fmt.Println("generating", compile)
	file, err = Root.Create(compile)
	if err != nil {
		return err
	}
	file.Close()
	Root.Chmod(compile, 0755)
	//fmt.Println(compile, "exists?", FileExists(detect))
	//fmt.Println(compile, "executable?", Executable(detect))

	return nil
}
