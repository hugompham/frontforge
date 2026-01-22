// Package preflight provides validation checks before project generation.
//
// Pre-flight checks verify that the environment is ready for project generation,
// preventing failures that would otherwise occur during file creation.
//
// Implemented checks:
//   - Node.js version (minimum v18.0.0)
//   - Package manager availability (npm, yarn, pnpm, bun)
//   - Directory conflict detection
//   - Disk space verification (minimum 500MB)
//
// Each check returns a CheckResult with pass/fail status, human-readable
// message, and optional suggestion for resolution. Checks can be marked
// as fatal, which prevents generation if they fail.
//
// Use RunAllChecks() to execute all validation checks and get aggregated results.
package preflight

import (
	"frontforge/internal/models"
)

// CheckResult represents the result of a single pre-flight check
type CheckResult struct {
	Name       string // Display name of the check (e.g., "Node.js Version")
	Passed     bool   // Whether the check passed
	Message    string // Detailed status message
	Suggestion string // Suggested action if check failed
	Fatal      bool   // If true, generation cannot proceed
}

// PreflightResults holds the results of all pre-flight checks
type PreflightResults struct {
	Checks     []CheckResult // All individual check results
	AllPassed  bool          // True if all checks passed
	FatalError bool          // True if any check failed with Fatal=true
}

// RunAllChecks executes all pre-flight validation checks
// and returns a consolidated result
func RunAllChecks(config models.Config) PreflightResults {
	var checks []CheckResult
	allPassed := true
	fatalError := false

	// Check 1: Node.js installation and version
	nodeCheck := CheckNodeJS()
	checks = append(checks, nodeCheck)
	if !nodeCheck.Passed {
		allPassed = false
		if nodeCheck.Fatal {
			fatalError = true
		}
	}

	// Check 2: Package manager availability
	pmCheck := CheckPackageManager(config.PackageManager)
	checks = append(checks, pmCheck)
	if !pmCheck.Passed {
		allPassed = false
		if pmCheck.Fatal {
			fatalError = true
		}
	}

	// Check 3: Directory conflicts
	dirCheck := CheckDirectoryConflicts(config.ProjectPath)
	checks = append(checks, dirCheck)
	if !dirCheck.Passed {
		allPassed = false
		if dirCheck.Fatal {
			fatalError = true
		}
	}

	// Check 4: Disk space
	diskCheck := CheckDiskSpace(config.ProjectPath)
	checks = append(checks, diskCheck)
	if !diskCheck.Passed {
		allPassed = false
		if diskCheck.Fatal {
			fatalError = true
		}
	}

	return PreflightResults{
		Checks:     checks,
		AllPassed:  allPassed,
		FatalError: fatalError,
	}
}
