package gologger

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Option func(logger *logrus.Logger) error

func WithDebugLevel(debug bool) Option {
	return func(logger *logrus.Logger) error {
		if debug {
			formatter := &logrus.TextFormatter{
				TimestampFormat: "02-01-2006 15:04:05",
				FullTimestamp:   true,
				ForceColors:     true,
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					return fmt.Sprintf("%s - ", formatFilePath(f.Function)), fmt.Sprintf(" %s:%d -", formatFilePath(f.File), f.Line)
				},
			}
			logger.SetLevel(logrus.DebugLevel)
			logger.SetFormatter(formatter)
			logger.SetReportCaller(true)
			logger.Level = logrus.DebugLevel
		}
		return nil
	}
}

func WithOutput(output io.Writer) Option {
	return func(logger *logrus.Logger) error {
		logger.SetOutput(output)
		return nil
	}
}

func WithFiles(outputFile string, errorFile string) Option {
	return func(logger *logrus.Logger) error {
		if _, err := os.Stat(outputFile); err == nil {
			os.Remove(outputFile)
		}
		if _, err := os.Stat(errorFile); err == nil {
			os.Remove(errorFile)
		}
		pathMap := lfshook.PathMap{
			logrus.ErrorLevel: errorFile,
			logrus.DebugLevel: outputFile,
			logrus.InfoLevel:  outputFile,
			logrus.WarnLevel:  outputFile,
		}
		logger.Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{
				TimestampFormat: "02-01-2006 15:04:05",
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					return fmt.Sprintf("%s ", formatFilePath(f.Function)), fmt.Sprintf(" %s:%d ", formatFilePath(f.File), f.Line)
				},
			},
		))
		return nil
	}
}

func WithFields(fields map[string]interface{}) Option {
	return func(logger *logrus.Logger) error {
		logger.WithFields(fields)
		return nil
	}
}

func WithField(msg string, val interface{}) Option {
	return func(logger *logrus.Logger) error {
		logEntry := logger.WithField(msg, val)
		logger = logEntry.Logger
		return nil
	}
}

func WithNullLogger() Option {
	return func(logger *logrus.Logger) error {
		logger.Out = io.Discard
		return nil
	}
}

func WithRuntimeContext() Option {
	return func(logger *logrus.Logger) error {
		if pc, file, line, ok := runtime.Caller(1); ok {
			fName := runtime.FuncForPC(pc).Name()
			logger.WithField("file", file).WithField("line", line).WithField("func", fName)
			return nil
		}
		return fmt.Errorf("logger option: failed to get runtime context")
	}
}
