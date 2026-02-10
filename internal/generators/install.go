package generators

import (
	"bufio"
	"fmt"
	"frontforge/internal/models"
	"io"
	"os/exec"
)

// RunInstall executes the package manager install command in the project directory
func RunInstall(projectPath string, config models.Config) error {
	var cmd *exec.Cmd

	// Determine the install command based on package manager
	switch config.PackageManager {
	case models.PackageManagerNpm:
		cmd = exec.Command("npm", "install")
	case models.PackageManagerYarn:
		cmd = exec.Command("yarn", "install")
	case models.PackageManagerPnpm:
		cmd = exec.Command("pnpm", "install")
	case models.PackageManagerBun:
		cmd = exec.Command("bun", "install")
	default:
		return fmt.Errorf("unsupported package manager: %s", config.PackageManager)
	}

	cmd.Dir = projectPath

	// Get stdout and stderr pipes for real-time streaming
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start install: %w", err)
	}

	// Stream output in real-time
	go streamOutput(stdout)
	go streamOutput(stderr)

	// Wait for completion
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("install failed: %w", err)
	}

	return nil
}

// streamOutput reads from a pipe and prints to stdout line by line
func streamOutput(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
