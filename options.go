package gologger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Option func(l *Logger) error

// calldepth is the call depth of the callsite function relative to the
// caller of the subsystem logger.  It is used to recover the filename and line
// number of the logging call if either the short or long file flags are
// specified.
const calldepth = 3

// WithDebugLevel sets the log level to debug with some preformatted output
func WithDebugLevel(debug bool) Option {
	return func(l *Logger) error {
		if debug {
			formatter := &logrus.TextFormatter{
				TimestampFormat: "02-01-2006 15:04:05",
				FullTimestamp:   true,
				ForceColors:     true,
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					return fmt.Sprintf("%s - ", formatFilePath(f.Function)), fmt.Sprintf(" %s:%d -", formatFilePath(f.File), f.Line)
				},
			}
			l.Logger.SetLevel(logrus.DebugLevel)
			l.Logger.SetFormatter(formatter)
			l.Logger.SetReportCaller(true)
		}
		return nil
	}
}

// WithLevel sets the log level and the corresponding format for the logger
func WithLevel(level Level) Option {
	return func(l *Logger) error {
		formatter := &logrus.TextFormatter{
			ForceColors:            true,
			PadLevelText:           false,
			DisableLevelTruncation: false,
		}
		switch level {
		case DebugLevel:
			formatter.TimestampFormat = "02-01-2006 15:04:05"
			formatter.FullTimestamp = true
			formatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
				return fmt.Sprintf("%s - ", formatFilePath(f.Function)), fmt.Sprintf(" %s:%d -", formatFilePath(f.File), f.Line)
			}
			l.Logger.SetReportCaller(true)
		case TraceLevel:
			formatter.TimestampFormat = "02-01-2006 15:04:05"
			formatter.FullTimestamp = true
			if pc, file, line, ok := runtime.Caller(calldepth); ok {
				fName := runtime.FuncForPC(pc).Name()
				l = l.WithField("file", file).WithField("line", line).WithField("func", fName)
			}
			formatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
				return fmt.Sprintf("%s - ", formatFilePath(f.Function)), fmt.Sprintf(" %s:%d -", formatFilePath(f.File), f.Line)
			}
			l.Logger.SetReportCaller(true)
		default:
			formatter.DisableTimestamp = true
			l.Logger.SetFormatter(formatter)
			l.Logger.SetReportCaller(false)

		}
		l.Logger.SetFormatter(formatter)
		l.Logger.SetLevel(logrus.Level(level))
		return nil
	}
}

// WithOutput sets the output for the logger
func WithOutput(output io.Writer) Option {
	return func(l *Logger) error {
		l.Logger.SetOutput(output)
		return nil
	}
}

// WithFiles configures the logger to write to the given files
func WithFiles(outputFile string, errorFile string) Option {
	return func(l *Logger) error {
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
		l.Logger.Hooks.Add(lfshook.NewHook(
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

// WithFields sets the fields for the logger
func WithFields(fields map[string]interface{}) Option {
	return func(l *Logger) error {
		l.Entry = l.WithFields(fields)
		return nil
	}
}

// WithField sets the field for the logger
func WithField(msg string, val interface{}) Option {
	return func(l *Logger) error {
		l = l.WithField(msg, val)
		return nil
	}
}

// WithNullLogger sets the logger to discard all output
func WithNullLogger() Option {
	return func(l *Logger) error {
		l.Logger.Out = io.Discard
		return nil
	}
}

// WithRuntimeContext sets the logger to include runtime context
func WithRuntimeContext() Option {
	return func(l *Logger) error {
		if pc, file, line, ok := runtime.Caller(1); ok {
			fName := runtime.FuncForPC(pc).Name()
			l = l.WithField("file", file).WithField("line", line).WithField("func", fName)
			formatter := &logrus.TextFormatter{
				TimestampFormat: "02-01-2006 15:04:05",
				FullTimestamp:   true,
				ForceColors:     true,
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					return fmt.Sprintf("%s - ", formatFilePath(f.Function)), fmt.Sprintf(" %s:%d -", formatFilePath(f.File), f.Line)
				},
			}
			l.Logger.SetFormatter(formatter)
			l.Logger.SetReportCaller(true)
			return nil
		}
		return fmt.Errorf("logger option: failed to get runtime context")
	}
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
