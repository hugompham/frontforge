package state

// LayoutState holds terminal dimensions and responsive layout settings
type LayoutState struct {
	// Terminal dimensions
	Width  int // Current terminal width
	Height int // Current terminal height

	// Adaptive layout
	AdaptiveBox int // Calculated box width based on terminal size
}

// NewLayoutState creates a LayoutState with safe defaults
func NewLayoutState() LayoutState {
	// Safe defaults for initial render before WindowSizeMsg
	initialWidth := 80
	initialHeight := 24

	return LayoutState{
		Width:       initialWidth,
		Height:      initialHeight,
		AdaptiveBox: calculateBoxWidth(initialWidth),
	}
}

// Update updates the layout state when terminal size changes
func (l *LayoutState) Update(width, height int) {
	l.Width = width
	l.Height = height
	l.AdaptiveBox = calculateBoxWidth(width)
}

// calculateBoxWidth determines adaptive box width based on terminal width
func calculateBoxWidth(terminalWidth int) int {
	const (
		minBoxWidth = 60
		maxBoxWidth = 100
		defaultBox  = 76
	)

	if terminalWidth < 80 {
		return minBoxWidth
	} else if terminalWidth > 120 {
		return maxBoxWidth
	}
	return defaultBox
}
