package logger

import (
	"app-skeleton-demo/internal/foundational/config"
	"log"
	"os"
	"strings"
)

// Logger provides a simple logging interface.
// In a real application, use a structured logger like slog or zap.
type Logger struct {
	*log.Logger
	level string
}

// New creates a new Logger instance.
func New(cfg config.LoggerConfig) *Logger {
	// This is a very simplified level handling for demonstration.
	// A real logger would have proper level parsing and filtering.
	l := log.New(os.Stdout, "[DemoApp] ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	l.Printf("FoundationalComponent: Logger initialized (simulated level: %s)\n", cfg.Level)
	return &Logger{Logger: l, level: strings.ToLower(cfg.Level)}
}

// Name returns the component name.
func (l *Logger) Name() string { return "Logger" }

// Debugf logs a debug message if the level is appropriate.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level == "debug" {
		l.Printf("DEBUG: "+format, v...)
	}
}

// Infof logs an info message.
func (l *Logger) Infof(format string, v ...interface{}) {
	// Assuming info is always printed for levels info and debug
	if l.level == "debug" || l.level == "info" {
		l.Printf("INFO: "+format, v...)
	}
}

// Errorf logs an error message.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Printf("ERROR: "+format, v...)
}

// Fatalf logs a fatal message and exits.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Fatalf("FATAL: "+format, v...)
}
