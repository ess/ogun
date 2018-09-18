package fs

import (
	"github.com/ess/conan"
)

func applicationPath(app conan.Application) string {
	return "/data/" + app.Name
}

func configPath(app conan.Application) string {
	return applicationPath(app) + "/shared/config"
}

func cachePath(app conan.Application) string {
	return applicationPath(app) + "/shared/build_cache"
}

func buildpackRoot() string {
	return "/engineyard/buildpacks/"
}

func buildpackPath(pack conan.Buildpack) string {
	return buildpackRoot() + pack.Name
}
