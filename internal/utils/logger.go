package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// InitLogger initializes the logger with configuration
func InitLogger(logLevel string) {
	Logger = logrus.New()
	
	// Set output to stdout
	Logger.SetOutput(os.Stdout)
	
	// Set formatter to JSON for structured logging
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	
	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Logger.SetLevel(logrus.InfoLevel)
		Logger.Warnf("Invalid log level '%s', defaulting to 'info'", logLevel)
	} else {
		Logger.SetLevel(level)
	}
}

// LogWithFields creates a log entry with fields
func LogWithFields(level logrus.Level, message string, fields map[string]interface{}) {
	entry := Logger.WithFields(logrus.Fields(fields))
	
	switch level {
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.FatalLevel:
		entry.Fatal(message)
	case logrus.PanicLevel:
		entry.Panic(message)
	default:
		entry.Info(message)
	}
}

