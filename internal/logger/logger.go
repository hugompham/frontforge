// Package logger provides simple structured logging for FrontForge CLI.
//
// The logger supports multiple log levels (Debug, Info, Warn, Error) and
// optional structured fields for context. By default, only Info and above
// are logged unless debug mode is enabled.
//
// Example usage:
//
//	log := logger.New(logger.LevelInfo, os.Stderr)
//	log.Info("Starting project generation", logger.F("project", "my-app"))
//	log.Error("Generation failed", logger.F("error", err))
package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a log message
type Level int

const (
	// LevelDebug shows detailed debugging information
	LevelDebug Level = iota
	// LevelInfo shows general informational messages
	LevelInfo
	// LevelWarn shows warning messages
	LevelWarn
	// LevelError shows error messages
	LevelError
)

// String returns the string representation of a log level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Field represents a structured key-value pair for logging context
type Field struct {
	Key   string
	Value interface{}
}

// F creates a new structured field for logging
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Logger provides structured logging with configurable output and levels
type Logger struct {
	level  Level
	output io.Writer
}

// New creates a new Logger with the specified minimum level and output writer
func New(level Level, output io.Writer) *Logger {
	if output == nil {
		output = os.Stderr
	}
	return &Logger{
		level:  level,
		output: output,
	}
}

// NewDefault creates a logger with Info level writing to stderr
func NewDefault() *Logger {
	return New(LevelInfo, os.Stderr)
}

// NewDebug creates a logger with Debug level writing to stderr
func NewDebug() *Logger {
	return New(LevelDebug, os.Stderr)
}

// SetLevel changes the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Debug logs a debug message with optional fields
func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

// Info logs an info message with optional fields
func (l *Logger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

// Warn logs a warning message with optional fields
func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

// Error logs an error message with optional fields
func (l *Logger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

// log handles the actual logging logic
func (l *Logger) log(level Level, msg string, fields ...Field) {
	// Skip if level is below threshold
	if level < l.level {
		return
	}

	// Format: [TIMESTAMP] LEVEL: message key=value key=value
	timestamp := time.Now().Format("15:04:05")
	output := fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), msg)

	// Add structured fields if present
	if len(fields) > 0 {
		for _, field := range fields {
			output += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	fmt.Fprintln(l.output, output)
}

// IsDebug returns true if debug level is enabled
func (l *Logger) IsDebug() bool {
	return l.level <= LevelDebug
}

// Global logger instance (can be replaced for testing or custom configuration)
var global = NewDefault()

// SetGlobal replaces the global logger instance
func SetGlobal(logger *Logger) {
	global = logger
}

// Global logger convenience functions

// Debug logs a debug message to the global logger
func Debug(msg string, fields ...Field) {
	global.Debug(msg, fields...)
}

// Info logs an info message to the global logger
func Info(msg string, fields ...Field) {
	global.Info(msg, fields...)
}

// Warn logs a warning message to the global logger
func Warn(msg string, fields ...Field) {
	global.Warn(msg, fields...)
}

// Error logs an error message to the global logger
func Error(msg string, fields ...Field) {
	global.Error(msg, fields...)
}

// SetLevel changes the global logger's minimum level
func SetLevel(level Level) {
	global.SetLevel(level)
}

// IsDebug returns true if the global logger has debug enabled
func IsDebug() bool {
	return global.IsDebug()
}
