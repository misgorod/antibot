package common

/*func NewLoggerMock() Logger {
	l := &loggerMock{
		DebugFunc:  func(args ...interface{}) {},
		TracefFunc: func(format string, args ...interface{}) {},
		ErrorFunc:  func(args ...interface{}) {},
		ErrorfFunc: func(format string, args ...interface{}) {},
	}
	l.WithErrorFunc = func(err error) Logger { return l }
	l.WithFieldFunc = func(key string, value interface{}) Logger { return l }
	return l
}

type loggerMock struct {
	DebugFunc     func(args ...interface{})
	TracefFunc    func(format string, args ...interface{})
	ErrorFunc     func(args ...interface{})
	ErrorfFunc    func(format string, args ...interface{})
	WithErrorFunc func(err error) Logger
	WithFieldFunc func(key string, value interface{}) Logger
}

func (l *loggerMock) Debug(args ...interface{}) {
	l.DebugFunc(args)
}

func (l *loggerMock) Tracef(format string, args ...interface{}) {
	l.TracefFunc(format, args...)
}

func (l *loggerMock) Error(args ...interface{}) {
	l.ErrorFunc(args...)
}

func (l *loggerMock) Errorf(format string, args ...interface{}) {
	l.ErrorfFunc(format, args...)
}

func (l *loggerMock) WithError(err error) Logger {
	return l.WithErrorFunc(err)
}
func (l *loggerMock) WithField(key string, value interface{}) Logger {
	return l.WithFieldFunc(key, value)
}
*/