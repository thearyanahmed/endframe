package shared

type LoggerInterface interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Trace(args ...interface{})
}
