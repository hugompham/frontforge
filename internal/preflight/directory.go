package preflight

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CheckDirectoryConflicts verifies the target directory is safe to use
func CheckDirectoryConflicts(projectPath string) CheckResult {
	result := CheckResult{
		Name:  "Directory Check",
		Fatal: false, // Not fatal - user can choose to continue with existing directory
	}

	// Check if path is empty
	if projectPath == "" {
		result.Passed = false
		result.Message = "Project path not specified"
		result.Suggestion = "Enter a project name or specify a path with -path flag"
		result.Fatal = true
		return result
	}

	// Get absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Invalid project path: %s", projectPath)
		result.Suggestion = "Use a valid directory path"
		result.Fatal = true
		return result
	}

	// Validate path safety - prevent system directory usage
	if err := ValidatePathSafety(absPath); err != nil {
		result.Passed = false
		result.Message = err.Error()
		result.Suggestion = "Choose a path in your home or project directory"
		result.Fatal = true
		return result
	}

	// Check if directory exists
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		// Directory doesn't exist - perfect for new project
		result.Passed = true
		result.Message = fmt.Sprintf("Target: %s (new directory) ✓", filepath.Base(absPath))
		result.Suggestion = ""
		return result
	}

	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Cannot access directory: %s", err.Error())
		result.Suggestion = "Check directory permissions"
		result.Fatal = true
		return result
	}

	// Path exists - check if it's a directory
	if !info.IsDir() {
		result.Passed = false
		result.Message = fmt.Sprintf("Path exists but is not a directory: %s", absPath)
		result.Suggestion = "Choose a different path or remove the existing file"
		result.Fatal = true
		return result
	}

	// Directory exists - check if it has files
	entries, err := os.ReadDir(absPath)
	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Cannot read directory: %s", err.Error())
		result.Suggestion = "Check directory permissions"
		result.Fatal = true
		return result
	}

	if len(entries) == 0 {
		// Empty directory - safe to use
		result.Passed = true
		result.Message = fmt.Sprintf("Target: %s (empty directory) ✓", filepath.Base(absPath))
		result.Suggestion = ""
		return result
	}

	// Directory exists and has files
	result.Passed = false
	result.Message = fmt.Sprintf("Directory already exists with %d file(s)", len(entries))
	result.Suggestion = "Choose a different directory name or remove existing files"
	return result
}

// validatePathSafety ensures the path is not in a dangerous system location
// ValidatePathSafety checks if a path is safe for project creation
func ValidatePathSafety(absPath string) error {
	normalizedPath := filepath.Clean(absPath)

	// Allow paths in system temporary directory
	tempDir := filepath.Clean(os.TempDir())
	if strings.HasPrefix(normalizedPath, tempDir+string(filepath.Separator)) ||
		normalizedPath == tempDir {
		return nil
	}

	// Forbidden system directories
	forbiddenPaths := []string{
		"/",
		"/bin",
		"/boot",
		"/dev",
		"/etc",
		"/lib",
		"/lib64",
		"/proc",
		"/root",
		"/sbin",
		"/sys",
		"/usr",
		"/var",
		"/System",           // macOS
		"/Library",          // macOS
		"/Applications",     // macOS
		"C:\\Windows",       // Windows
		"C:\\Program Files", // Windows
	}

	for _, forbidden := range forbiddenPaths {
		normalizedForbidden := filepath.Clean(forbidden)
		if normalizedPath == normalizedForbidden ||
			strings.HasPrefix(normalizedPath, normalizedForbidden+string(filepath.Separator)) {
			return fmt.Errorf("cannot create project in system directory: %s", forbidden)
		}
	}

	return nil
}
