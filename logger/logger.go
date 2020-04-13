package logger

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

var Logger *logrus.Logger

func init() {
	Logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.InfoLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}

	// logging level
	lvl, _ := os.LookupEnv("LOG_LEVEL")

	//if !ok {
	//	lvl = "debug"
	//}
	if lvl == "" {
		lvl = "info"
	}
	SetLevel(lvl)

	// logging json formatter
	fmt, _ := os.LookupEnv("LOG_FMT")

	if fmt == "json" {
		AsJson()
	}
}

func SetLevel(level string) {
	var parsedLevel logrus.Level
	var err error
	if parsedLevel, err = logrus.ParseLevel(level); err != nil {
		parsedLevel = logrus.DebugLevel
	}
	Logger.SetLevel(parsedLevel)
}

func AsJson() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
}

// Debug logs a message at debug level.
func Debugf(format string, v ...interface{}) {
	Logger.Debugf(format, v...)
}

// Info logs a message at info level.
func Infof(format string, v ...interface{}) {
	Logger.Infof(format, v...)
}

// Warn logs a message at warn level.
func Warnf(format string, v ...interface{}) {
	Logger.Warnf(format, v...)
}

// Fatal logs a message at fatal level.
// Fatalf is the equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	Logger.Fatalf(format, v...)
}

// Error logs a message at fatal level.
// Errorf is the equivalent to Printf() followed by a call to os.Exit(1).
func Errorf(format string, v ...interface{}) {
	Logger.Errorf(format, v...)
}
