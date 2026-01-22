package generators

import (
	"fmt"
	"frontforge/internal/errors"
	"os"
	"path/filepath"
	"strings"
)

// NormalizePath converts a user-provided path to an absolute, clean path
// Handles relative paths, ".." segments, and validates safety
func NormalizePath(userPath, cwd string) (string, error) {
	if userPath == "" {
		return "", errors.NewPathError("", "path cannot be empty", nil)
	}

	var absPath string
	var err error

	// Handle absolute vs relative paths
	if filepath.IsAbs(userPath) {
		absPath = userPath
	} else {
		// Resolve relative to current working directory
		absPath = filepath.Join(cwd, userPath)
	}

	// Clean the path (removes "..", ".", duplicate slashes)
	absPath = filepath.Clean(absPath)

	// Ensure path is absolute after cleaning
	if !filepath.IsAbs(absPath) {
		absPath, err = filepath.Abs(absPath)
		if err != nil {
			return "", errors.NewPathError(userPath, "failed to resolve absolute path", err)
		}
	}

	// Validate the path is safe (not in system directories)
	if err := validatePathSafety(absPath); err != nil {
		return "", err
	}

	return absPath, nil
}

// ValidateProjectPath checks if a path is suitable for project creation
// Returns an error if the path is unsafe or already contains files
func ValidateProjectPath(absPath string) error {
	// Check if path is safe
	if err := validatePathSafety(absPath); err != nil {
		return err
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		// Path doesn't exist - this is fine, we'll create it
		return nil
	}

	if err != nil {
		return errors.NewPathError(absPath, "cannot access path", err)
	}

	// Path exists - check if it's a directory
	if !info.IsDir() {
		return errors.NewPathError(absPath, "path exists but is not a directory", nil)
	}

	// Directory exists - check if it's empty
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return errors.NewPathError(absPath, "cannot read directory", err)
	}

	if len(entries) > 0 {
		return errors.NewPathError(absPath,
			fmt.Sprintf("directory is not empty (%d file(s) found)", len(entries)), nil)
	}

	return nil
}

// IsPathSafe checks if a path is safe for project creation
// Returns true if the path is safe, false if it's in a system directory
func IsPathSafe(absPath string) bool {
	return validatePathSafety(absPath) == nil
}

// validatePathSafety ensures the path is not in a dangerous location
func validatePathSafety(absPath string) error {
	// Normalize path for comparison
	normalizedPath := filepath.Clean(absPath)

	// Allow paths in system temporary directory (e.g., /var/folders on macOS)
	tempDir := filepath.Clean(os.TempDir())
	if strings.HasPrefix(normalizedPath, tempDir+string(filepath.Separator)) ||
		normalizedPath == tempDir {
		return nil
	}

	// List of forbidden paths (system directories)
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

	// Check if path starts with any forbidden path
	for _, forbidden := range forbiddenPaths {
		// Normalize forbidden path for comparison
		normalizedForbidden := filepath.Clean(forbidden)

		// Check if the path is exactly the forbidden path or a subdirectory
		if normalizedPath == normalizedForbidden ||
			strings.HasPrefix(normalizedPath, normalizedForbidden+string(filepath.Separator)) {
			return errors.NewPathError(absPath,
				fmt.Sprintf("cannot create project in system directory: %s", forbidden), nil)
		}
	}

	return nil
}

// GetProjectName extracts the project name from a path
// Returns the last component of the path (directory name)
func GetProjectName(absPath string) string {
	return filepath.Base(absPath)
}

// EnsureParentDir ensures the parent directory of a path exists
// Creates parent directories if they don't exist
func EnsureParentDir(filePath string) error {
	parentDir := filepath.Dir(filePath)

	// Check if parent directory exists
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		// Create parent directory with appropriate permissions
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return errors.NewPathError(filePath,
				"failed to create parent directory", err)
		}
	}

	return nil
}
