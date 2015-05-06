package kasi

import (
	"os"

	"github.com/op/go-logging"
)

var logFormat = logging.MustStringFormatter(
	"%{color}%{time:2006-01-01T15:04:05.000} %{shortfunc} > [%{level:.2s}]%{color:reset} %{message}",
)

var log *logging.Logger
var DefaultLogLevel logging.Level = logging.CRITICAL

func SetLogging(level logging.Level) *logging.Logger {
	log = logging.MustGetLogger("kasi")

	logging.SetFormatter(logFormat)

	log_backend := logging.NewLogBackend(os.Stdout, "", 0)
	log_backend.Color = true

	log_backend_level := logging.AddModuleLevel(log_backend)
	log_backend_level.SetLevel(level, "")

	log.SetBackend(log_backend_level)

	return log
}
