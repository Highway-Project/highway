package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

//default setting
var Logger = &logrus.Logger{}

func InitLogger(logLevel string, prettyLog bool) {
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
	})
	Logger.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)
}
