package logger

import (
	"os"
	
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	
	// Set default configuration
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false,
	})
	
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

// SetLevel sets the logging level
func SetLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

// SetFormatter sets the log formatter
func SetFormatter(format string) {
	switch format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			PrettyPrint:     false,
		})
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	default:
		log.SetFormatter(&logrus.JSONFormatter{})
	}
}

// WithFields creates a log entry with fields
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return log.WithFields(fields)
}

// Debug logs a debug message
func Debug(msg string) {
	log.Debug(msg)
}

// Info logs an info message
func Info(msg string) {
	log.Info(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	log.Warn(msg)
}

// Error logs an error message
func Error(msg string) {
	log.Error(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string) {
	log.Fatal(msg)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
