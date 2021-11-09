package beans

var (
	defaultLogger = &loggingHandler{}
	logger        ILogger
)

// SetLogger sets an implementation of ILogger to be used as the logger for the
// beans package
func SetLogger(l ILogger) {
	logger = l
}

// LogCallbacks returns a handler that allows to register individual callbacks
// to be used by the beans package to report errors, info messages and/or debug
// messages
func LogCallbacks() ILogCallback {
	return defaultLogger
}

type loggingHandler struct {
	onErrHandler   ErrorCallback
	onInfoHandler  MessageCallback
	onDebugHandler MessageCallback
}

func log() ILogger {
	if logger != nil {
		return logger
	}
	return defaultLogger
}

func (l *loggingHandler) SetErrorCallback(callback ErrorCallback) {
	l.onErrHandler = callback
}

func (l *loggingHandler) SetInfoCallback(callback MessageCallback) {
	l.onInfoHandler = callback
}

func (l *loggingHandler) SetDebugCallback(callback MessageCallback) {
	l.onDebugHandler = callback
}

func (l *loggingHandler) Error(err error) {
	if l.onErrHandler != nil {
		l.onErrHandler(err)
	}
}

func (l *loggingHandler) Info(msg string) {
	if l.onInfoHandler != nil {
		l.onInfoHandler(msg)
	}
}

func (l *loggingHandler) Debug(msg string) {
	if l.onDebugHandler != nil {
		l.onDebugHandler(msg)
	}
}
