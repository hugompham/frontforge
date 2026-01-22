// Package errors provides structured error types for FrontForge.
//
// This package defines three main error types:
//   - GenerationError: Errors during project file generation
//   - PathError: Errors related to path validation and file system operations
//   - PreflightError: Errors from pre-flight validation checks
//
// All error types implement the error interface and support error unwrapping
// for compatibility with errors.Unwrap() and errors.Is().
package errors

import (
	"fmt"
)

// GenerationError represents errors that occur during project generation
type GenerationError struct {
	Stage   string // Which stage failed (e.g., "directory_creation", "file_generation")
	Message string // Human-readable error message
	Cause   error  // Underlying error if any
}

// Error implements the error interface
func (e *GenerationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("generation failed at %s: %s (cause: %v)", e.Stage, e.Message, e.Cause)
	}
	return fmt.Sprintf("generation failed at %s: %s", e.Stage, e.Message)
}

// Unwrap returns the underlying error for error wrapping chains
func (e *GenerationError) Unwrap() error {
	return e.Cause
}

// NewGenerationError creates a new GenerationError
func NewGenerationError(stage, message string, cause error) *GenerationError {
	return &GenerationError{
		Stage:   stage,
		Message: message,
		Cause:   cause,
	}
}

// PreflightError represents errors that occur during pre-flight validation
type PreflightError struct {
	Check      string // Which check failed (e.g., "node_version", "disk_space")
	Message    string // Human-readable error message
	Suggestion string // Suggested action to resolve the error
	Fatal      bool   // Whether this error prevents generation
}

// Error implements the error interface
func (e *PreflightError) Error() string {
	if e.Fatal {
		return fmt.Sprintf("preflight check failed [FATAL]: %s - %s", e.Check, e.Message)
	}
	return fmt.Sprintf("preflight check failed [WARNING]: %s - %s", e.Check, e.Message)
}

// NewPreflightError creates a new PreflightError
func NewPreflightError(check, message, suggestion string, fatal bool) *PreflightError {
	return &PreflightError{
		Check:      check,
		Message:    message,
		Suggestion: suggestion,
		Fatal:      fatal,
	}
}

// PathError represents errors related to path validation and manipulation
type PathError struct {
	Path    string // The problematic path
	Message string // Human-readable error message
	Cause   error  // Underlying error if any
}

// Error implements the error interface
func (e *PathError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("path error for '%s': %s (cause: %v)", e.Path, e.Message, e.Cause)
	}
	return fmt.Sprintf("path error for '%s': %s", e.Path, e.Message)
}

// Unwrap returns the underlying error for error wrapping chains
func (e *PathError) Unwrap() error {
	return e.Cause
}

// NewPathError creates a new PathError
func NewPathError(path, message string, cause error) *PathError {
	return &PathError{
		Path:    path,
		Message: message,
		Cause:   cause,
	}
}
