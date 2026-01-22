//go:build unix || darwin || linux
// +build unix darwin linux

package preflight

import (
	"fmt"
	"path/filepath"
	"syscall"
)

const (
	// MinDiskSpaceBytes is the minimum required disk space (500 MB)
	MinDiskSpaceBytes = 500 * 1024 * 1024
	// RecommendedDiskSpaceBytes is the recommended disk space (1 GB)
	RecommendedDiskSpaceBytes = 1024 * 1024 * 1024
)

// CheckDiskSpace verifies sufficient disk space is available
func CheckDiskSpace(projectPath string) CheckResult {
	result := CheckResult{
		Name:  "Disk Space",
		Fatal: true, // Insufficient disk space will cause generation to fail
	}

	// If path is empty, use current directory
	checkPath := projectPath
	if checkPath == "" {
		checkPath = "."
	}

	// Get absolute path
	absPath, err := filepath.Abs(checkPath)
	if err != nil {
		result.Passed = false
		result.Message = "Cannot determine disk space"
		result.Suggestion = "Check path validity"
		return result
	}

	// Get directory to check (use parent if target doesn't exist yet)
	statPath := absPath
	var statInfo syscall.Stat_t
	if err := syscall.Stat(absPath, &statInfo); err != nil {
		// Path doesn't exist, check parent directory
		statPath = filepath.Dir(absPath)
	}

	// Get filesystem stats
	var stat syscall.Statfs_t
	err = syscall.Statfs(statPath, &stat)
	if err != nil {
		result.Passed = false
		result.Message = fmt.Sprintf("Cannot check disk space: %s", err.Error())
		result.Suggestion = "Verify path permissions"
		return result
	}

	// Calculate available space
	// stat.Bavail is available blocks for non-root users
	// stat.Bsize is block size in bytes
	availableBytes := stat.Bavail * uint64(stat.Bsize)

	// Convert to human-readable format
	availableMB := float64(availableBytes) / (1024 * 1024)
	availableGB := availableMB / 1024

	// Check if sufficient space is available
	if availableBytes < MinDiskSpaceBytes {
		result.Passed = false
		result.Message = fmt.Sprintf("Insufficient disk space: %.1f MB available (requires 500 MB)",
			availableMB)
		result.Suggestion = "Free up disk space before continuing"
		return result
	}

	// Warn if below recommended space
	if availableBytes < RecommendedDiskSpaceBytes {
		result.Passed = true
		result.Message = fmt.Sprintf("Limited disk space: %.1f MB available (%.1f GB recommended)",
			availableMB, float64(RecommendedDiskSpaceBytes)/(1024*1024*1024))
		result.Suggestion = "Consider freeing up more space for optimal performance"
		result.Fatal = false
		return result
	}

	// Sufficient space available
	result.Passed = true
	if availableGB >= 1.0 {
		result.Message = fmt.Sprintf("Disk space: %.1f GB available ✓", availableGB)
	} else {
		result.Message = fmt.Sprintf("Disk space: %.0f MB available ✓", availableMB)
	}
	result.Suggestion = ""
	return result
}
