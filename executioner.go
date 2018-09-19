package ogun

type Executioner interface {
	Execute(string) ([]byte, error)
}
