package support

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Level logrus.Level
}

var logLevelMap = map[string]logrus.Level{
	"trace": logrus.TraceLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

func GetLogger(args *Arguments) *logrus.Logger {
	var logger = logrus.New()
	level, ok := logLevelMap[args.LogLevel]
	if !ok {
		level = logLevelMap["info"]
	}
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.WithFields(logrus.Fields{
		"args": args,
	}).Info("Given options")
	return logger
}
