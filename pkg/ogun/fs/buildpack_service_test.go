package fs

import (
	//"fmt"
	"testing"

	"github.com/ess/ogun/pkg/ogun"
)

func setupDetect() {
	logger.Reset()
	setupFs()
	stubApplication()
}

func TestBuildpackService_Detect(t *testing.T) {
	service := NewBuildpackService(logger)

	runner.Reset()
	service.runner = runner

	t.Run("when there's a custom buildpack", func(t *testing.T) {
		setupDetect()

		pack := ogun.Buildpack{Name: "custom", Location: applicationPath(app) + "/shared/cached_copy/.ogun/buildpack"}

		err := stubPack(pack)
		if err != nil {
			t.Errorf("couldn't stub the pack")
		}

		result, err := service.Detect(app)

		t.Run("it returns the app's custom buildpack", func(t *testing.T) {
			if result.Location != pack.Location {
				t.Errorf("expected %s, got %s", pack.Location, result.Location)
			}
		})

		t.Run("it returns no error", func(t *testing.T) {
			if err != nil {
				t.Errorf("expected no error, got %s", err)
			}
		})

	})

	t.Run("when there's no custom buildpack", func(t *testing.T) {
		t.Run("but there are no buildpacks installed", func(t *testing.T) {
			setupDetect()

			_, err := service.Detect(app)

			t.Run("it returns an error", func(t *testing.T) {
				if err == nil {
					t.Errorf("expected an error")
				}
			})
		})

		t.Run("and there are buildpacks installed", func(t *testing.T) {
			pack1 := ogun.Buildpack{Name: "buildpack-pack1", Location: "/engineyard/buildpacks/buildpack-pack1"}
			pack2 := ogun.Buildpack{Name: "buildpack-pack2", Location: "/engineyard/buildpacks/buildpack-pack2"}

			setupDetect()
			err := stubPack(pack1)
			if err != nil {
				t.Errorf("could not stub pack1")
			}

			err = stubPack(pack2)
			if err != nil {
				t.Errorf("could not stub pack2")
			}

			t.Run("but none of them can build the app", func(t *testing.T) {
				runner.Reset()

				_, derr := service.Detect(app)

				t.Run("it returns an error", func(t *testing.T) {
					if derr == nil {
						t.Errorf("expected an error")
					}
				})
			})

			t.Run("but more than one can build the app", func(t *testing.T) {
				runner.Reset()
				runner.Add(pack1.Location + "/bin/detect")
				runner.Add(pack2.Location + "/bin/detect")

				_, derr := service.Detect(app)

				t.Run("it returns an error", func(t *testing.T) {
					if derr == nil {
						t.Errorf("expected an error")
					}
				})
			})

			t.Run("and exactly one can build the app", func(t *testing.T) {
				runner.Reset()
				runner.Add(pack1.Location + "/bin/detect")

				result, derr := service.Detect(app)

				t.Run("it returns the correct buildpack", func(t *testing.T) {
					if result.Location != pack1.Location {
						t.Errorf("expected %s, got %s", pack1.Location, result.Location)
					}
				})

				t.Run("it returns no error", func(t *testing.T) {
					if derr != nil {
						t.Errorf("expected no error, got %s", derr)
					}
				})
			})
		})
	})

}
