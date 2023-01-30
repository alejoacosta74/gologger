// package log implements a custom logger based on logrus
package gologger

import (
	"io"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

var (
	loggerOnce  sync.Once
	logger      *Logger
	loggerError error
)

// Returns true if the log level is set to 'debug' or 'trace'
func (l *Logger) IsDebug() bool {
	return l.GetLevel() >= DebugLevel
}

// Returns the log level
func (l *Logger) GetLevel() Level {
	level := l.Logger.GetLevel()
	return Level(level)
}

// Sets the log level
func (l *Logger) SetLevel(level Level) {
	l.Logger.SetLevel(logrus.Level(level))
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.Entry = l.Entry.WithField(key, value)
	return l
}

// Same as WithField but implements the ILogger interface
func (l *Logger) With(key string, value interface{}) ILogger {
	return l.WithField(key, value)
}

// Returns a new logger instance with the given field
func (l *Logger) NewLoggerWithField(key string, value interface{}) *Logger {
	newLogger := l.Entry.WithField(key, value)
	return &Logger{newLogger}
}

// Returns a new logger instantiated from the existing logger with the given fields
func (l *Logger) NewLoggerWithFields(fields map[string]interface{}) *Logger {
	newLogger := l.Entry.WithFields(fields)
	return &Logger{newLogger}
}

// Returns a singleton logger instance with the given options
func NewLogger(opts ...Option) (*Logger, error) {
	loggerOnce.Do(func() {
		logger, loggerError = createNewLogger(opts...)
	})
	return logger, loggerError
}

func createNewLogger(opts ...Option) (*Logger, error) {
	l := logrus.New()
	formatter := &logrus.TextFormatter{
		DisableTimestamp: true,
		ForceColors:      true,
	}
	l.SetFormatter(formatter)

	// by default set output to no-op
	l.SetOutput(io.Discard)

	logger := &Logger{
		Entry: logrus.NewEntry(l),
	}

	for _, opt := range opts {
		if err := opt(logger); err != nil {
			return nil, err
		}
	}
	return logger, nil
}
