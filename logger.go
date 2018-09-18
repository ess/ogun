package conan

import (
	"io"
)

type Logger interface {
	Info(string, string)
	Error(string, string)
	AddOutput(io.Writer)
	Writers() []io.Writer
}
