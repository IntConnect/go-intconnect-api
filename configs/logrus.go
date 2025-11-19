package configs

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// Singleton instance
var (
	loggerInstance *logrus.Logger
	once           sync.Once // Ensure initialization only once
)

// GetLogger returns singleton logger instance
func GetLogger() *logrus.Logger {
	once.Do(func() {
		loggerInstance = initLogger()
	})
	return loggerInstance
}

// Initialize logger (called only once)
func initLogger() *logrus.Logger {
	log := logrus.New()

	// Determine environment
	isDev := os.Getenv("APP_ENV") != "production"

	if isDev {
		// Development configuration
		log.SetFormatter(&logrus.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		log.SetOutput(os.Stdout)
		log.SetLevel(logrus.DebugLevel)
	} else {
		// Production configuration
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05Z07:00",
		})

	}

	return log
}

// Convenience methods (optional - untuk kemudahan akses)

func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}

// WithFields helper
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return GetLogger().WithFields(logrus.Fields(fields))
}

// WithError helper
func WithError(err error) *logrus.Entry {
	return GetLogger().WithError(err)
}

// Deprecated: Use GetLogger() directly if you need access to the logger instance
func Logger() *logrus.Logger {
	return GetLogger()
}
