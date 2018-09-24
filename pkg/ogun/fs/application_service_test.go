package fs

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/ess/testscope"
	"github.com/spf13/afero"
)

func TestApplicationService_Get(t *testing.T) {
	testscope.SkipUnlessUnit(t)
	appBase := "/data"
	appName := gofakeit.Generate("????????????")
	service := &ApplicationService{}

	t.Run("when the application does not exist", func(t *testing.T) {
		Root = afero.NewMemMapFs()

		_, getErr := service.Get(appName)

		t.Run("it returns an error", func(t *testing.T) {
			if getErr == nil {
				t.Errorf("expected an error")
			}
		})
	})

	t.Run("when the application exists", func(t *testing.T) {
		Root = afero.NewMemMapFs()

		err := CreateDir(appBase+"/"+appName, 0755)
		if err != nil {
			t.Errorf("eould not set up the app directory")
		}

		result, getErr := service.Get(appName)

		t.Run("it returns a populated Application", func(t *testing.T) {
			if result.Name != appName {
				t.Errorf("expected name '%s', got '%s'", appName, result.Name)
			}
		})

		t.Run("it returns no error", func(t *testing.T) {
			if getErr != nil {
				t.Errorf("expected no error, got %s", getErr)
			}
		})

	})

}
