package log

import (
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/log"
)

const (
	EnvDevelopment = "DEVELOPMENT"
	EnvProduction  = "PRODUCTION"
)

// Package scoped Logger variable
var l *Logger
var o sync.Once

// Logger is a logger that wraps a Logger
type Logger struct {
	*log.Logger
}

// NewLog returns a new Logger singleton instance based on the provided application environment
func NewLog(appEnv string) *Logger {
	initLogger(appEnv, os.Stdout)
	return l
}

// NewLogWithFile returns a new Logger singleton instance based on the provided application environment and log file path
func NewLogWithFile(appEnv string, logFile *os.File) *Logger {
	initLogger(appEnv, logFile)
	return l
}

// Debug logs a debug message
func Debug(message interface{}) {
	l.Logger.Debug(message)
}

// Info logs an info message
func Info(message interface{}) {
	l.Logger.Info(message)
}

// Warn logs a warning message
func Warn(message interface{}) {
	l.Logger.Warn(message)
}

// Error logs an error message
func Error(message interface{}) {
	l.Logger.Error(message)
}

// Fatal logs an fatal message and exit the application
func Fatal(message interface{}) {
	l.Logger.Fatal(message)
}

// initLogger initializes the logger based on the provided application environment and writer
func initLogger(appEnv string, w io.Writer) {
	if l == nil {
		o.Do(func() {
			logger := log.NewWithOptions(w, log.Options{
				ReportTimestamp: true,
				Prefix:          "mpwt 🍊 ",
			})
			switch appEnv {
			case EnvDevelopment:
				logger.SetLevel(log.DebugLevel)
				logger.SetReportCaller(true)
			case EnvProduction:
				logger.SetLevel(log.InfoLevel)
			default:
				logger.SetLevel(log.InfoLevel)
			}
			l = &Logger{logger}
		})
	}
}
