package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Log *logrus.Entry
)

func Init(logLevel string) {
	logEntry := logrus.NewEntry(logrus.New())

	var level logrus.Level
	switch logLevel {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	}
	logEntry.Logger.SetLevel(level)
	logEntry.Logger.SetOutput(os.Stdout)
	logEntry.Logger.SetFormatter(&logrus.JSONFormatter{})
	logEntry.Logger.SetReportCaller(true)
	Log = logEntry
}
