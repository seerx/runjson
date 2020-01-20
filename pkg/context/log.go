package context

type Log interface {
	Info(format string, val ...interface{})
	Warn(format string, val ...interface{})
	Error(err error, format string, val ...interface{})
	Debug(format string, val ...interface{})
}
