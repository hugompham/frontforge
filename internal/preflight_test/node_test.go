package preflight_test

import (
	"frontforge/internal/preflight"
	"testing"
)

func TestCheckNodeJS(t *testing.T) {
	result := preflight.CheckNodeJS()

	// Should either pass or fail with clear message
	if result.Name == "" {
		t.Error("CheckResult should have a name")
	}

	if result.Message == "" {
		t.Error("CheckResult should have a message")
	}

	// Fatal should be true for Node.js check
	if result.Fatal != true {
		t.Error("Node.js check should be marked as fatal")
	}

	// If failed, should have suggestion
	if !result.Passed && result.Suggestion == "" {
		t.Error("Failed check should provide a suggestion")
	}
}
