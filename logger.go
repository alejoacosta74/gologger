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

const DebugLevel = logrus.DebugLevel

var (
	loggerOnce  sync.Once
	logger      *Logger
	loggerError error
)

func (l *Logger) IsDebug() bool {
	return l.GetLevel() == logrus.DebugLevel
}

func (l *Logger) GetLevel() logrus.Level {
	return l.Logger.GetLevel()
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.Entry = l.Entry.WithField(key, value)
	return l
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.Entry = l.Entry.WithFields(fields)
	return l
}

func (l *Logger) SetFields(fields map[string]interface{}) {
	l.Entry = l.Entry.WithFields(fields)
}

func Fields(f map[string]interface{}) logrus.Fields {
	return f

}

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
