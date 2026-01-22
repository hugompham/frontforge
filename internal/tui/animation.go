package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
)

// ============================================================================
// FORGING ANIMATION FRAMES - Custom spinner for the forging process
// ============================================================================

// ForgingSpinner returns a simple progress spinner
func ForgingSpinner() spinner.Spinner {
	return spinner.Spinner{
		Frames: []string{
			"⠋",
			"⠙",
			"⠹",
			"⠸",
			"⠼",
			"⠴",
			"⠦",
			"⠧",
			"⠇",
			"⠏",
		},
		FPS: 4, // Slow spinner for minimal distraction
	}
}

// ============================================================================
// ANIMATED TEXT HELPERS - For simple feedback animations
// ============================================================================

// AnimatedDots returns an animated "..." with cycling visibility
// tickCount: current tick counter
// maxDots: maximum number of dots (default 3)
// This provides subtle feedback without being distracting (<300ms updates)
func AnimatedDots(tickCount int, maxDots int) string {
	if maxDots == 0 {
		maxDots = 3
	}

	// Update every 60 ticks (~1 second at 60fps) for minimal distraction
	numDots := (tickCount / AnimatedDotsInterval) % (maxDots + 1)
	dots := ""
	for i := 0; i < numDots; i++ {
		dots += "."
	}

	// Pad to maintain width consistency
	for i := numDots; i < maxDots; i++ {
		dots += " "
	}

	return dots
}

// WorkingIndicator returns a rotating work-in-progress indicator
// Provides simple visual feedback that something is happening
func WorkingIndicator(tickCount int) string {
	indicators := []string{"[|]", "[/]", "[-]", "[\\]"}
	return indicators[(tickCount/WorkingIndicatorSpeed)%len(indicators)]
}
