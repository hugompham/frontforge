package meta

import "fmt"

// ScaffoldError represents a failure when running an upstream CLI.
type ScaffoldError struct {
	Framework string
	Command   string
	ExitCode  int
	Stderr    string // last 50 lines
}

func (e *ScaffoldError) Error() string {
	if e.Command == "" {
		return fmt.Sprintf("scaffold error for %s: %s", e.Framework, e.Stderr)
	}
	return fmt.Sprintf("scaffold error for %s (exit %d): %s\nCommand: %s",
		e.Framework, e.ExitCode, e.Stderr, e.Command)
}
