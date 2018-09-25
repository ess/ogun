package ogun

type Release struct {
	Name        string
	Application Application
	Buildpack   Buildpack
}

type ReleaseService interface {
	Create(string, Application) (Release, error)
	Build(Release, Buildpack) error
	Clean(Release) error
	Package(Release) error
	Delete(Release) error
}
