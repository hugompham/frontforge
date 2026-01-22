package logger_test

import (
	"bytes"
	"frontforge/internal/logger"
	"strings"
	"testing"
)

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(logger.LevelInfo, &buf)

	// Debug should not be logged at Info level
	log.Debug("debug message")
	if buf.String() != "" {
		t.Errorf("debug message should not be logged at Info level")
	}

	// Info should be logged
	log.Info("info message")
	if !strings.Contains(buf.String(), "INFO") || !strings.Contains(buf.String(), "info message") {
		t.Errorf("info message should be logged: %s", buf.String())
	}

	buf.Reset()

	// Warn should be logged
	log.Warn("warning message")
	if !strings.Contains(buf.String(), "WARN") || !strings.Contains(buf.String(), "warning message") {
		t.Errorf("warning message should be logged: %s", buf.String())
	}

	buf.Reset()

	// Error should be logged
	log.Error("error message")
	if !strings.Contains(buf.String(), "ERROR") || !strings.Contains(buf.String(), "error message") {
		t.Errorf("error message should be logged: %s", buf.String())
	}
}

func TestLoggerDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(logger.LevelDebug, &buf)

	// Debug should be logged at Debug level
	log.Debug("debug message")
	if !strings.Contains(buf.String(), "DEBUG") || !strings.Contains(buf.String(), "debug message") {
		t.Errorf("debug message should be logged at Debug level: %s", buf.String())
	}
}

func TestLoggerFields(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(logger.LevelInfo, &buf)

	// Log with structured fields
	log.Info("operation complete",
		logger.F("project", "my-app"),
		logger.F("duration", "5s"),
	)

	output := buf.String()
	if !strings.Contains(output, "project=my-app") {
		t.Errorf("output should contain project field: %s", output)
	}
	if !strings.Contains(output, "duration=5s") {
		t.Errorf("output should contain duration field: %s", output)
	}
}

func TestLoggerSetLevel(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(logger.LevelInfo, &buf)

	// Initially Debug should not log
	log.Debug("debug1")
	if buf.String() != "" {
		t.Errorf("debug should not log initially")
	}

	// Change to Debug level
	log.SetLevel(logger.LevelDebug)

	// Now Debug should log
	log.Debug("debug2")
	if !strings.Contains(buf.String(), "debug2") {
		t.Errorf("debug should log after SetLevel: %s", buf.String())
	}
}

func TestIsDebug(t *testing.T) {
	infoLog := logger.New(logger.LevelInfo, nil)
	debugLog := logger.New(logger.LevelDebug, nil)

	if infoLog.IsDebug() {
		t.Error("Info logger should not have debug enabled")
	}

	if !debugLog.IsDebug() {
		t.Error("Debug logger should have debug enabled")
	}
}

func TestGlobalLogger(t *testing.T) {
	var buf bytes.Buffer
	customLog := logger.New(logger.LevelDebug, &buf)
	logger.SetGlobal(customLog)

	// Global functions should use the custom logger
	logger.Info("test message")

	if !strings.Contains(buf.String(), "test message") {
		t.Errorf("global logger should use custom instance: %s", buf.String())
	}
}
