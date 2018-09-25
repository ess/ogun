package ogun

type Buildpack struct {
	Name     string
	Location string
}

type BuildpackService interface {
	Detect(Application) (Buildpack, error)
}
