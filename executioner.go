package conan

type Executioner interface {
	Execute(string) ([]byte, error)
}
