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
	Fatalf(message string)
	Errorf(message string)
	Warnf(message string)
	Infof(message string)
	Debugf(message string)
	Tracef(message string)
}
