package meta

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ExecScaffold runs an external command with a timeout, capturing output.
// If dryRun is true, returns the command that would have been executed without running it.
func ExecScaffold(framework string, dryRun bool, name string, args ...string) error {
	cmdStr := fmt.Sprintf("%s %s", name, strings.Join(args, " "))

	if dryRun {
		fmt.Printf("[meta-scaffold] Would run: %s\n", cmdStr)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}

		// Take last 50 lines of stderr
		stderrStr := stderr.String()
		lines := strings.Split(stderrStr, "\n")
		if len(lines) > 50 {
			lines = lines[len(lines)-50:]
		}

		return &ScaffoldError{
			Framework: framework,
			Command:   cmdStr,
			ExitCode:  exitCode,
			Stderr:    strings.Join(lines, "\n"),
		}
	}

	return nil
}

// ExecInDir runs a command in a specific directory with timeout.
func ExecInDir(dir, framework string, dryRun bool, name string, args ...string) error {
	cmdStr := fmt.Sprintf("%s %s", name, strings.Join(args, " "))

	if dryRun {
		fmt.Printf("[meta-scaffold] Would run: %s (in %s)\n", cmdStr, dir)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}

		stderrStr := stderr.String()
		lines := strings.Split(stderrStr, "\n")
		if len(lines) > 50 {
			lines = lines[len(lines)-50:]
		}

		return &ScaffoldError{
			Framework: framework,
			Command:   cmdStr,
			ExitCode:  exitCode,
			Stderr:    strings.Join(lines, "\n"),
		}
	}

	return nil
}
