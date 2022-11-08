package logger

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type loggerConfig interface {
	GetLogLevel() string
}

// Setup initialize the Logrus logger configuration.
func Setup(conf loggerConfig) {
	log.SetOutput(os.Stdout)
	log.SetLevel(getLogLevelFromString(conf.GetLogLevel()))
}

func Logger() *log.Logger {
	return &log.Logger{}
}

func getLogLevelFromString(level string) log.Level {
	switch strings.ToLower(level) {
	case "fatal":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	case "trace":
		return log.TraceLevel
	default:
		return log.PanicLevel
	}
}
