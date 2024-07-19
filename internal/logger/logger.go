package logger

import "github.com/sirupsen/logrus"

var Log *logrus.Logger

func InitializeLogger() {
	Log = logrus.New()

	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)
}
