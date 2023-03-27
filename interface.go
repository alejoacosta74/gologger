package gologger

type ILogger interface {
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
	Tracef(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	IsDebug() bool
	With(key string, value interface{}) ILogger
	SetLevel(level Level)
}

// interface assertion
var _ ILogger = &Logger{}
