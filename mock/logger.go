package mock

import (
	"fmt"
	"io"
)

type Entry struct {
	Severity string
	Context  string
	Message  string
}

func (e *Entry) String() string {
	return fmt.Sprintf("[%s](%s) - %s", e.Severity, e.Context, e.Message)
}

type Logger struct {
	Entries []*Entry
}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) log(severity string, context string, message string) {
	e := &Entry{Severity: severity, Context: context, Message: message}

	logger.Entries = append(logger.Entries, e)
}

func (logger *Logger) Info(context string, message string) {
	logger.log("INFO", context, message)
}

func (logger *Logger) Error(context string, message string) {
	logger.log("ERROR", context, message)
}

func (logger *Logger) AddOutput(output io.Writer) {}

func (logger *Logger) Writers() []io.Writer {
	return nil
}

func (logger *Logger) Reset() {
	logger.Entries = make([]*Entry, 0)
}
