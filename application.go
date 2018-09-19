package ogun

type Application struct {
	Name    string
	BaseDir string
}

func (app Application) String() string {
	return app.Name
}

type ApplicationService interface {
	Get(string) (Application, error)
}
