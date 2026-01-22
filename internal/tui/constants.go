package tui

// UI timing and animation constants
const (
	// Animation update intervals (in ticks)
	AnimatedDotsInterval  = 60 // Ticks between dot updates in loading states
	WorkingIndicatorSpeed = 8  // Ticks between rotation of working indicator

	// Display timing
	MinForgingDisplayTime = 5 // Minimum seconds to display forging animation

	// Terminal dimensions
	MinTerminalWidth = 80 // Minimum recommended terminal width

	// File system permissions
	DirPermissions  = 0755 // Standard directory permissions (rwxr-xr-x)
	FilePermissions = 0644 // Standard file permissions (rw-r--r--)
)

// Progress tracking constants
const (
	// Question catalog size - total number of possible form questions
	// This is used for progress tracking in the blueprint phase
	MaxQuestionCount = 18 // Maximum number of questions in custom mode
)

// Terminal width calculation
const (
	// Adaptive box width thresholds
	NarrowTerminalThreshold = 80  // Below this, use minimum box width
	WideTerminalThreshold   = 120 // Above this, use maximum box width
)
