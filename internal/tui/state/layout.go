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
	// AdaptiveBox will be calculated by parent package using CalculateBoxWidth
	return LayoutState{
		Width:       80,
		Height:      24,
		AdaptiveBox: 76, // Default box width
	}
}

// Update updates the layout state when terminal size changes
// Note: AdaptiveBox should be calculated by parent using CalculateBoxWidth
func (l *LayoutState) Update(width, height int) {
	l.Width = width
	l.Height = height
	// AdaptiveBox is set separately by parent package
}
