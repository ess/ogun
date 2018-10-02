package fs

import (
	//"fmt"
	"testing"

	"github.com/ess/ogun/pkg/ogun"
)

func setupCreate() {
	logger.Reset()
	setupFs()
	stubApplication()
	setupRelease()
}

func setupRelease() {
	release = ogun.Release{Name: releaseName, Application: app}

	if DirectoryExists(releasePath(release)) {
		Root.RemoveAll(releasePath(release))
	}
}

func TestReleaseService_Create(t *testing.T) {
	service := NewReleaseService(logger)

	t.Run("it creates the release path", func(t *testing.T) {
		setupCreate()

		if DirectoryExists(releasePath(release)) {
			t.Errorf("expected %s not to exist yet", releasePath(release))
		}

		service.Create(releaseName, app)

		if !DirectoryExists(releasePath(release)) {
			t.Errorf("expected %s to exist after creation", releasePath(release))
		}
	})

	t.Run("it copies the source cache to the release path", func(t *testing.T) {
		setupCreate()

		ohai := releasePath(release) + "/ohai"
		if FileExists(ohai) {
			t.Errorf("expected %s not to exist yet", ohai)
		}

		service.Create(releaseName, app)

		if !FileExists(ohai) {
			t.Errorf("expected %s to exist after creation", ohai)
		}
	})

	t.Run("it returns the release", func(t *testing.T) {
		setupCreate()

		result, _ := service.Create(releaseName, app)

		if result.Name != releaseName {
			t.Errorf("expected release %s, got release %s", releaseName, result.Name)
		}
	})
}
