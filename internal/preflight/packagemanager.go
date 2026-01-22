package preflight

import (
	"fmt"
	"os/exec"
	"strings"
)

// CheckPackageManager verifies that the selected package manager is available
func CheckPackageManager(packageManager string) CheckResult {
	result := CheckResult{
		Name:  fmt.Sprintf("Package Manager (%s)", packageManager),
		Fatal: true, // Package manager is required for dependency installation
	}

	// Determine the command name based on package manager choice
	var cmdName string
	switch packageManager {
	case "npm":
		cmdName = "npm"
	case "yarn":
		cmdName = "yarn"
	case "pnpm":
		cmdName = "pnpm"
	case "bun":
		cmdName = "bun"
	default:
		result.Passed = false
		result.Message = fmt.Sprintf("Unknown package manager: %s", packageManager)
		result.Suggestion = "Select a valid package manager (npm, yarn, pnpm, or bun)"
		return result
	}

	// Check if command exists
	cmdPath, err := exec.LookPath(cmdName)
	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("%s not found", cmdName)
		result.Suggestion = getPackageManagerInstallSuggestion(packageManager)
		return result
	}

	// Get version to verify it's working
	cmd := exec.Command(cmdPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to check %s version", cmdName)
		result.Suggestion = fmt.Sprintf("Verify %s installation", cmdName)
		return result
	}

	version := strings.TrimSpace(string(output))
	result.Passed = true
	result.Message = fmt.Sprintf("%s v%s ✓", cmdName, version)
	result.Suggestion = ""
	return result
}

// getPackageManagerInstallSuggestion provides installation instructions for each package manager
func getPackageManagerInstallSuggestion(packageManager string) string {
	switch packageManager {
	case "npm":
		return "npm is bundled with Node.js. Install Node.js from https://nodejs.org"
	case "yarn":
		return "Install yarn: npm install -g yarn or visit https://yarnpkg.com"
	case "pnpm":
		return "Install pnpm: npm install -g pnpm or visit https://pnpm.io"
	case "bun":
		return "Install bun: curl -fsSL https://bun.sh/install | bash or visit https://bun.sh"
	default:
		return "Select a different package manager or install the selected one"
	}
}
