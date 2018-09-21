package ogun

type Executioner interface {
	Execute(string, []Variable) ([]byte, error)
}
