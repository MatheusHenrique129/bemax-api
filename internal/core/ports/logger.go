package ports

type Logger interface {
	Info(message string, tags ...interface{})
	Warn(message string, tags ...interface{})
	Debug(message string, args ...interface{})
	Fatal(message string, args ...interface{})
	Error(message string, err error, tags ...interface{})
}
