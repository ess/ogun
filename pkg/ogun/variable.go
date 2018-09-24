package ogun

type Variable string

type VariableService interface {
	All(Application) []Variable
}
