package kennel

type Suite interface {
	Step(interface{}, interface{})
	BeforeScenario(func(interface{}))
	BeforeSuite(func())
	AfterSuite(func())
}

type Stepper interface {
	StepUp(Suite)
}

var steppers []Stepper

func Register(s Stepper) {
	steppers = append(steppers, s)
}

func StepUp(s Suite) {
	for _, stepper := range steppers {
		stepper.StepUp(s)
	}
}
