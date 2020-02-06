package fw

type LogLevel int

const (
	LogFatal LogLevel = iota
	LogError
	LogWarn
	LogInfo
	LogDebug
	LogTrace
	LogOff
)

type LogLevelName string

const (
	LogFatalName LogLevelName = "Fatal"
	LogErrorName LogLevelName = "Error"
	LogWarnName  LogLevelName = "Warn"
	LogInfoName  LogLevelName = "Info"
	LogDebugName LogLevelName = "Debug"
	LogTraceName LogLevelName = "Trace"
)

type Logger interface {
	Fatal(message string)
	Error(err error)
	Warn(message string)
	Info(message string)
	Debug(message string)
	Trace(message string)
	Fatalf(message string, a ...interface{})
	Errorf(message string, a ...interface{})
	Warnf(message string, a ...interface{})
	Infof(message string, a ...interface{})
	Debugf(message string, a ...interface{})
	Tracef(message string, a ...interface{})
}
