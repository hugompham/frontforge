package preflight_test

import (
	"frontforge/internal/models"
	"frontforge/internal/preflight"
	"frontforge/internal/testutil"
	"testing"
)

func TestRunAllChecks(t *testing.T) {
	tempDir := testutil.TempDir(t)

	config := models.Config{
		ProjectPath:    tempDir,
		PackageManager: "npm",
	}

	results := preflight.RunAllChecks(config)

	// Should run all checks
	if len(results.Checks) == 0 {
		t.Fatal("Expected some checks to run")
	}

	// Should have at least 4 checks (Node, PM, Directory, Disk)
	if len(results.Checks) < 4 {
		t.Errorf("Expected at least 4 checks, got %d", len(results.Checks))
	}

	// AllPassed should match individual results
	expectedAllPassed := true
	for _, check := range results.Checks {
		if !check.Passed {
			expectedAllPassed = false
			break
		}
	}

	if results.AllPassed != expectedAllPassed {
		t.Errorf("AllPassed=%v but individual checks suggest %v", results.AllPassed, expectedAllPassed)
	}

	// FatalError should be set if any fatal check failed
	expectedFatalError := false
	for _, check := range results.Checks {
		if !check.Passed && check.Fatal {
			expectedFatalError = true
			break
		}
	}

	if results.FatalError != expectedFatalError {
		t.Errorf("FatalError=%v but checks suggest %v", results.FatalError, expectedFatalError)
	}
}
