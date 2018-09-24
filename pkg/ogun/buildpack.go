package ogun

type Buildpack struct {
	Name     string
	Location string
}

type BuildpackService interface {
	Get(string) (Buildpack, error)
	Detect(Application) (Buildpack, error)
}
