//go:build windows
// +build windows

package preflight

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
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
	if _, err := os.Stat(absPath); err != nil {
		// Path doesn't exist, check parent directory
		statPath = filepath.Dir(absPath)
	}

	// Windows requires the path to end with a backslash for root directories
	// For GetDiskFreeSpaceEx, we use the volume root
	volumePath := filepath.VolumeName(statPath) + "\\"

	// Get disk space using Windows API
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")

	var freeBytesAvailable uint64
	var totalBytes uint64
	var totalFreeBytes uint64

	volumePathPtr, err := syscall.UTF16PtrFromString(volumePath)
	if err != nil {
		result.Passed = false
		result.Message = "Cannot check disk space"
		result.Suggestion = "Verify path validity"
		return result
	}

	ret, _, _ := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(volumePathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)

	if ret == 0 {
		result.Passed = false
		result.Message = "Cannot check disk space"
		result.Suggestion = "Verify path permissions"
		return result
	}

	availableBytes := freeBytesAvailable

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
		result.Message = fmt.Sprintf("Disk space: %.1f GB available", availableGB)
	} else {
		result.Message = fmt.Sprintf("Disk space: %.0f MB available", availableMB)
	}
	result.Suggestion = ""
	return result
}
