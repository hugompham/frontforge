package meta

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// ProbeUpstreamCLI attempts to detect the version of an upstream CLI tool.
// Returns the version string, or "" if the tool is not found or errors.
func ProbeUpstreamCLI(command string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
