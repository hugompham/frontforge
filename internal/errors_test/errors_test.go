package errors_test

import (
	stderr "errors"
	"frontforge/internal/errors"
	"frontforge/internal/testutil"
	"strings"
	"testing"
)

func TestGenerationError(t *testing.T) {
	baseErr := stderr.New("file not found")
	genErr := errors.NewGenerationError("package.json", "failed to write file", baseErr)

	// Should have descriptive error message
	errMsg := genErr.Error()
	if !strings.Contains(errMsg, "package.json") {
		t.Errorf("error message should contain stage: %s", errMsg)
	}
	if !strings.Contains(errMsg, "failed to write file") {
		t.Errorf("error message should contain message: %s", errMsg)
	}

	// Should unwrap to original error
	unwrapped := stderr.Unwrap(genErr)
	if unwrapped != baseErr {
		t.Errorf("expected unwrap to return base error, got %v", unwrapped)
	}
}

func TestPathError(t *testing.T) {
	baseErr := stderr.New("permission denied")
	pathErr := errors.NewPathError("/tmp/test", "cannot access directory", baseErr)

	// Should have descriptive error message
	errMsg := pathErr.Error()
	if !strings.Contains(errMsg, "/tmp/test") {
		t.Errorf("error message should contain path: %s", errMsg)
	}
	if !strings.Contains(errMsg, "cannot access directory") {
		t.Errorf("error message should contain message: %s", errMsg)
	}

	// Should unwrap to original error
	unwrapped := stderr.Unwrap(pathErr)
	if unwrapped != baseErr {
		t.Errorf("expected unwrap to return base error, got %v", unwrapped)
	}
}

func TestPreflightError(t *testing.T) {
	preflightErr := errors.NewPreflightError("Node.js", "Node.js is required but not installed", "Install Node.js v18 or later", true)

	// Should have descriptive error message
	errMsg := preflightErr.Error()
	if !strings.Contains(errMsg, "Node.js") {
		t.Errorf("error message should contain check name: %s", errMsg)
	}
	if !strings.Contains(errMsg, "Node.js is required but not installed") {
		t.Errorf("error message should contain message: %s", errMsg)
	}

	// Fatal error should contain FATAL marker
	if !strings.Contains(errMsg, "FATAL") {
		t.Errorf("fatal error message should contain FATAL marker: %s", errMsg)
	}

	// Non-fatal error should contain WARNING marker
	warningErr := errors.NewPreflightError("PackageManager", "npm is recommended", "Install npm", false)
	if !strings.Contains(warningErr.Error(), "WARNING") {
		t.Errorf("warning error message should contain WARNING marker: %s", warningErr.Error())
	}
}

func TestErrorWithoutCause(t *testing.T) {
	// Errors without a cause should work fine
	genErr := errors.NewGenerationError("config", "invalid configuration", nil)
	testutil.AssertError(t, genErr)

	errMsg := genErr.Error()
	if !strings.Contains(errMsg, "config") || !strings.Contains(errMsg, "invalid configuration") {
		t.Errorf("error message missing expected content: %s", errMsg)
	}

	// Unwrap should return nil
	unwrapped := stderr.Unwrap(genErr)
	if unwrapped != nil {
		t.Errorf("expected unwrap to return nil for error without cause, got %v", unwrapped)
	}
}
