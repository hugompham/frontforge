package preflight

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	MinNodeMajorVersion = 20
	MinNodeMinorVersion = 19
	MinNodePatchVersion = 0
)

// CheckNodeJS verifies that Node.js is installed and meets minimum version requirements
func CheckNodeJS() CheckResult {
	result := CheckResult{
		Name:  "Node.js Installation",
		Fatal: true, // Node.js is required for project generation
	}

	// Check if node command exists
	nodePath, err := exec.LookPath("node")
	if err != nil {
		result.Passed = false
		result.Message = "Node.js not found"
		result.Suggestion = "Install Node.js v20.19.0 or higher from https://nodejs.org"
		return result
	}

	// Get Node.js version
	cmd := exec.Command(nodePath, "--version")
	output, err := cmd.Output()
	if err != nil {
		result.Passed = false
		result.Message = "Failed to check Node.js version"
		result.Suggestion = "Verify Node.js installation"
		return result
	}

	version := strings.TrimSpace(string(output))

	// Parse version (format: v18.19.0)
	major, minor, patch, err := parseNodeVersion(version)
	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Invalid Node.js version format: %s", version)
		result.Suggestion = "Reinstall Node.js from https://nodejs.org"
		return result
	}

	// Check if version meets minimum requirements
	if !meetsMinimumVersion(major, minor, patch) {
		result.Passed = false
		result.Message = fmt.Sprintf("Node.js %s found (requires v%d.%d.%d+)",
			version, MinNodeMajorVersion, MinNodeMinorVersion, MinNodePatchVersion)
		result.Suggestion = fmt.Sprintf("Update Node.js to v%d.%d.%d or higher",
			MinNodeMajorVersion, MinNodeMinorVersion, MinNodePatchVersion)
		return result
	}

	result.Passed = true
	result.Message = fmt.Sprintf("Node.js %s ✓", version)
	result.Suggestion = ""
	return result
}

// parseNodeVersion extracts major, minor, patch from version string
// Expected format: v18.19.0
func parseNodeVersion(version string) (major, minor, patch int, err error) {
	// Remove 'v' prefix if present
	version = strings.TrimPrefix(version, "v")

	// Match semantic version pattern
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(version)

	if len(matches) != 4 {
		return 0, 0, 0, fmt.Errorf("invalid version format: %s", version)
	}

	major, err = strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, 0, err
	}

	minor, err = strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, 0, err
	}

	patch, err = strconv.Atoi(matches[3])
	if err != nil {
		return 0, 0, 0, err
	}

	return major, minor, patch, nil
}

// meetsMinimumVersion checks if the given version meets minimum requirements
func meetsMinimumVersion(major, minor, patch int) bool {
	if major > MinNodeMajorVersion {
		return true
	}
	if major < MinNodeMajorVersion {
		return false
	}

	// Major version matches, check minor
	if minor > MinNodeMinorVersion {
		return true
	}
	if minor < MinNodeMinorVersion {
		return false
	}

	// Major and minor match, check patch
	return patch >= MinNodePatchVersion
}
