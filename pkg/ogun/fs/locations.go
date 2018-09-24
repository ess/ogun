package fs

import (
	"github.com/ess/ogun/pkg/ogun"
)

func applicationPath(app ogun.Application) string {
	return "/data/" + app.Name
}

func configPath(app ogun.Application) string {
	return applicationPath(app) + "/shared/config"
}

func cachePath(app ogun.Application) string {
	return applicationPath(app) + "/shared/build_cache"
}

func buildpackRoot() string {
	return "/engineyard/buildpacks/"
}

func buildpackPath(pack ogun.Buildpack) string {
	return buildpackRoot() + pack.Name
}
