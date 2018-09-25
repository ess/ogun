package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log     *logrus.Logger
	writers []io.Writer
}

func NewLogger() *Logger {
	log := logrus.New()
	log.Formatter = &Formatter{}

	logger := &Logger{log: log, writers: make([]io.Writer, 0)}
	logger.AddOutput(os.Stdout)
	return logger
}

func (logger *Logger) Info(context string, message string) {
	logger.log.WithField("context", context).Info(message)
}

func (logger *Logger) Error(context string, message string) {
	logger.log.WithField("context", context).Error(message)
}

func (logger *Logger) AddOutput(output io.Writer) {
	logger.writers = append(logger.writers, output)
	logger.log.Out = io.MultiWriter(logger.writers...)
}

func (logger *Logger) Writers() []io.Writer {
	writers := make([]io.Writer, 0)

	for _, writer := range logger.writers {
		writers = append(writers, writer)
	}

	return writers
}
