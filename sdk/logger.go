package sdk

type Logger interface {
	Info(pattern string, args ...interface{})
	Error(pattern string, args ...interface{})
	Debug(pattern string, args ...interface{})
}

type DefaultLogger struct {
	logger Logger
}

var defaultLogger = &DefaultLogger{}

func SetLogger(logger Logger) {
	defaultLogger.logger = logger
}

func (l *DefaultLogger) Info(pattern string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Info(pattern, args...)
	}
}

func (l *DefaultLogger) Error(pattern string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Error(pattern, args...)
	}
}

func (l *DefaultLogger) Debug(pattern string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Debug(pattern, args...)
	}
}
