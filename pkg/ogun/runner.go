package ogun

type Runner interface {
	Execute(string, []Variable) ([]byte, error)
}
