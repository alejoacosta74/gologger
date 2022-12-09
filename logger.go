// package log implements a custom logger based on logrus
package gologger

import (
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

type Level uint32

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

var (
	loggerOnce  sync.Once
	logger      *Logger
	loggerError error
)

// Returns true if the log level is set to debug
func (l *Logger) IsDebug() bool {
	return l.GetLevel() == logrus.DebugLevel
}

// Returns the log level
func (l *Logger) GetLevel() logrus.Level {
	return l.Logger.GetLevel()
}

// Returns a new logger instance with the given field
func (l *Logger) NewLoggerWithField(key string, value interface{}) *Logger {
	newLogger := l.Entry.WithField(key, value)
	return &Logger{newLogger}
}

// Returns a new logger instance with the given fields
func (l *Logger) NewLoggerWithFields(fields map[string]interface{}) *Logger {
	newLogger := l.Entry.WithFields(fields)
	return &Logger{newLogger}
}

// Returns a new logger instance with the given options
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

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
