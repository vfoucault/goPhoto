package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func SetupLog() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stderr)

	// logging level
	lvl, ok := os.LookupEnv("LOGLEVEL")

	if !ok || lvl == "" {
		logrus.SetLevel(logrus.InfoLevel)
	}
	if lvl == "DEBUG" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if lvl == "ERROR" {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	// LogFormat
	format, _ := os.LookupEnv("LOGFORMAT")
	if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		})
	}
}
