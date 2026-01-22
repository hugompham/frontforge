package preflight

import (
	"fmt"
	"os"
	"path/filepath"
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
