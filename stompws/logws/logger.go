package logws

import (
	"log"
	"os"
)

// Logger interface for logging
type Logger interface {
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// DefaultLogger implements Logger using standard log package
type DefaultLogger struct {
	logger *log.Logger
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger() Logger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "[STOMP] ", log.LstdFlags),
	}
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf("[INFO] "+format, args...)
}

func (l *DefaultLogger) Warningf(format string, args ...interface{}) {
	l.logger.Printf("[WARN] "+format, args...)
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+format, args...)
}

// NoOpLogger implements Logger with no-op operations
type NoOpLogger struct{}

func (l *NoOpLogger) Infof(format string, args ...interface{})    {}
func (l *NoOpLogger) Warningf(format string, args ...interface{}) {}
func (l *NoOpLogger) Errorf(format string, args ...interface{})   {}
