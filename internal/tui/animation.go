package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
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
		FPS: 4, // Much slower spinner - was 10
	}
}

// HeatingSpinner returns a spinner for initial heating/loading phase
// Animation depicts rising heat and flames
func HeatingSpinner() spinner.Spinner {
	return spinner.Spinner{
		Frames: []string{
			"[~]                    ",
			"   [~]                 ",
			"      [~]              ",
			"         [~]           ",
			"            [~]        ",
			"         [~]           ",
			"      [~]              ",
			"   [~]                 ",
		},
		FPS: 8, // Faster for heating effect
	}
}

// TemperingSpinner returns a spinner for cooling/tempering phase
// Animation depicts cooling down process
func TemperingSpinner() spinner.Spinner {
	return spinner.Spinner{
		Frames: []string{
			"[v]  [v]  [v]  [v]     ", // Drips cooling
			" [v]  [v]  [v]  [v]    ",
			"  [v]  [v]  [v]  [v]   ",
			"   [v]  [v]  [v]  [v]  ",
			"    [v]  [v]  [v]  [v] ",
			"     [v]  [v]  [v]  [v]",
		},
		FPS: 3, // Slower for calming effect
	}
}

// HammerSpinner returns a simple hammering animation
// Simpler than ForgingSpinner for less intensive moments
func HammerSpinner() spinner.Spinner {
	return spinner.Spinner{
		Frames: []string{
			"[H]  ",
			" [H] ",
			"  [H]",
			" [H] ",
		},
		FPS: 4,
	}
}

// SparksSpinner returns a sparks-only animation
// For quick transitions or light activity
func SparksSpinner() spinner.Spinner {
	return spinner.Spinner{
		Frames: []string{
			"[*]     ",
			" [*]    ",
			"  [*]   ",
			"   [*]  ",
			"    [*] ",
			"   [*]  ",
			"  [*]   ",
			" [*]    ",
		},
		FPS: 10, // Fast for sparking effect
	}
}

// ============================================================================
// PROGRESS ANIMATION FRAMES - For visual feedback during operations
// ============================================================================

// GetForgeAnimationFrame returns the current frame for the forging animation
// tickCount: current tick counter from the model
// Returns: Simple progress indicator
func GetForgeAnimationFrame(tickCount int) string {
	frames := []string{
		"●○○○○",
		"○●○○○",
		"○○●○○",
		"○○○●○",
		"○○○○●",
		"○○○●○",
		"○○●○○",
		"○●○○○",
	}

	// Cycle through frames
	frameIndex := (tickCount / 10) % len(frames)
	return frames[frameIndex]
}

// GetHeatingAnimationFrame returns the current frame for heating animation
func GetHeatingAnimationFrame(tickCount int) string {
	frames := []string{
		"       [~]                 ",
		"          [~]              ",
		"             [~]           ",
		"                [~]        ",
		"             [~]           ",
		"          [~]              ",
	}

	frameIndex := (tickCount / 8) % len(frames)
	return frames[frameIndex]
}

// GetForgingAnimationFrame returns ASCII art of hammer striking anvil
// Uses Harmonica spring physics for smooth hammer motion
// intensity: Current hammer position from spring physics (0.0 to 1.0)
func GetForgingAnimationFrame(intensity float64) string {
	// Choose frame based on intensity (hammer position)
	var hammerFrame int
	if intensity < 0.3 {
		hammerFrame = 0 // Hammer raised high
	} else if intensity < 0.5 {
		hammerFrame = 1 // Hammer mid-swing
	} else if intensity < 0.7 {
		hammerFrame = 2 // Hammer near anvil
	} else {
		hammerFrame = 3 // Hammer striking (with sparks)
	}

	frames := []string{
		// Frame 0: Hammer raised high
		`
       [H]
        |
        |
    ════════════
     [ANVIL]`,
		// Frame 1: Hammer mid-swing
		`

       [H]
        |
    ════════════
     [ANVIL]`,
		// Frame 2: Hammer near anvil
		`


       [H]
    ════════════
     [ANVIL]`,
		// Frame 3: Hammer striking (with sparks)
		`


      *[H]*
    ════════════
     [ANVIL]`,
	}

	return frames[hammerFrame]
}

// GetWelcomeAnimationFrame returns ASCII art of a forge furnace with animated flames
// Uses Harmonica spring physics for smooth, natural breathing fire effect
// intensity: Current flame intensity from spring physics (0.0 to 1.0)
func GetWelcomeAnimationFrame(intensity float64) string {
	// Choose frame based on intensity from spring physics
	var flameLevel int
	if intensity < 0.4 {
		flameLevel = 0 // Low flames
	} else if intensity < 0.6 {
		flameLevel = 1 // Medium-low flames
	} else if intensity < 0.8 {
		flameLevel = 2 // Medium-high flames
	} else {
		flameLevel = 3 // High flames
	}

	frames := []string{
		// Frame 0: Low flames - all frames MUST have same height (7 lines)
		`

    ╔══════════╗
    ║   ▒▒▒    ║
    ║  ▒▒▒▒▒   ║
    ║ ▓▓▓▓▓▓▓  ║
    ╚══════════╝`,
		// Frame 1: Medium-low flames
		`
       ▒▒
    ╔══════════╗
    ║  ▒▒▒▒▒   ║
    ║ ▒▒▒▒▒▒▒  ║
    ║ ▓▓▓▓▓▓▓  ║
    ╚══════════╝`,
		// Frame 2: Medium-high flames
		`
      ▒▒▒▒
    ╔══════════╗
    ║ ▒▒▒▒▒▒▒  ║
    ║ ▒▒▒▒▒▒▒  ║
    ║ ▓▓▓▓▓▓▓  ║
    ╚══════════╝`,
		// Frame 3: High flames
		`     ▒▒▒▒▒▒
      ▒▒▒▒
    ╔══════════╗
    ║ ▒▒▒▒▒▒▒  ║
    ║ ▒▒▒▒▒▒▒  ║
    ║ ▓▓▓▓▓▓▓  ║
    ╚══════════╝`,
	}

	return frames[flameLevel]
}

// ============================================================================
// ANIMATED TEXT HELPERS - For dynamic text effects
// ============================================================================

// AnimatedDots returns an animated "..." with cycling visibility
// tickCount: current tick counter
// maxDots: maximum number of dots (default 3)
func AnimatedDots(tickCount int, maxDots int) string {
	if maxDots == 0 {
		maxDots = 3
	}

	// Much slower animation - update every 60 ticks (~1 second at 60fps) instead of every 10 ticks
	numDots := (tickCount / 60) % (maxDots + 1)
	dots := ""
	for i := 0; i < numDots; i++ {
		dots += "."
	}

	// Pad to maintain width
	for i := numDots; i < maxDots; i++ {
		dots += " "
	}

	return dots
}

// PulsingText returns text that pulses in intensity
// Cycles between normal and bold styling
func PulsingText(text string, tickCount int, normalStyle, boldStyle lipgloss.Style) string {
	// Pulse every 20 ticks (~320ms)
	if (tickCount/20)%2 == 0 {
		return boldStyle.Render(text)
	}
	return normalStyle.Render(text)
}

// ============================================================================
// PROGRESS BAR ANIMATIONS - Enhanced progress indicators
// ============================================================================

// AnimatedProgressBar returns a progress bar with animated filling effect
// current: current progress value
// total: total progress value
// width: width of the progress bar in characters
// tickCount: current tick for animation
func AnimatedProgressBar(current, total, width, tickCount int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(float64(width) * percentage)
	if filled > width {
		filled = width
	}
	empty := width - filled

	// Animated filling character that cycles
	filledChars := []string{"█", "▓", "▒"}
	fillingChar := filledChars[(tickCount/5)%len(filledChars)]

	// Build the bar
	bar := ""
	for i := 0; i < filled-1; i++ {
		bar += "█"
	}

	// Last filled character is animated if not complete
	if filled > 0 && percentage < 1.0 {
		bar += fillingChar
	} else if filled > 0 {
		bar += "█"
	}

	// Empty portion
	for i := 0; i < empty; i++ {
		bar += "░"
	}

	return bar
}

// ============================================================================
// STATUS INDICATORS - Animated status symbols
// ============================================================================

// WorkingIndicator returns a rotating work-in-progress indicator
func WorkingIndicator(tickCount int) string {
	indicators := []string{"[|]", "[/]", "[-]", "[\\]"}
	return indicators[(tickCount/8)%len(indicators)]
}

// LoadingIndicator returns a loading animation
func LoadingIndicator(tickCount int) string {
	indicators := []string{"[.]  ", "[..]", "[...]", "[..]"}
	return indicators[(tickCount/10)%len(indicators)]
}

// SuccessAnimation returns frames for success animation
// Plays once when shown
func SuccessAnimation(frame int) string {
	frames := []string{
		"                          ", // Frame 0: empty
		"         [*]              ", // Frame 1: single spark
		"      [*] [*] [*]         ", // Frame 2: sparks spread
		"   [*] [*] [+] [*] [*]    ", // Frame 3: checkmark appears
		"      [✓] SUCCESS          ", // Frame 4: final state
	}

	if frame >= len(frames) {
		return frames[len(frames)-1]
	}
	return frames[frame]
}

// ErrorAnimation returns frames for error animation
func ErrorAnimation(frame int) string {
	frames := []string{
		"                          ",
		"         [!]              ",
		"      [!] [X] [!]         ",
		"   [!] [X] [✗] [X] [!]    ",
		"      [✗] FAILED           ",
	}

	if frame >= len(frames) {
		return frames[len(frames)-1]
	}
	return frames[frame]
}

// ============================================================================
// SPRING ANIMATIONS - Smooth, natural motion using Harmonica
// ============================================================================

// SpringAnimator handles spring-based smooth animations
// Provides natural, physics-based motion for value transitions
type SpringAnimator struct {
	spring   harmonica.Spring
	value    float64
	velocity float64
	target   float64
}

// NewSpringAnimator creates a new spring animator with default settings
// Parameters optimized for TUI animations (<300ms feel, critically damped)
func NewSpringAnimator() *SpringAnimator {
	return &SpringAnimator{
		// FPS(60): Update at 60 FPS for smooth motion
		// 8.0: Angular velocity - moderately fast (higher = faster)
		// 1.0: Damping ratio - critically damped (no oscillation, fastest without overshoot)
		spring:   harmonica.NewSpring(harmonica.FPS(60), 8.0, 1.0),
		value:    0,
		velocity: 0,
		target:   0,
	}
}

// NewSpringAnimatorCustom creates a spring animator with custom physics
// angularVelocity: Higher = faster animation (try 6.0 - 12.0)
// dampingRatio: <1 = springy/bouncy, 1 = smooth, >1 = slow
func NewSpringAnimatorCustom(angularVelocity, dampingRatio float64) *SpringAnimator {
	return &SpringAnimator{
		spring:   harmonica.NewSpring(harmonica.FPS(60), angularVelocity, dampingRatio),
		value:    0,
		velocity: 0,
		target:   0,
	}
}

// SetTarget sets a new target value to animate towards
func (sa *SpringAnimator) SetTarget(target float64) {
	sa.target = target
}

// Update advances the animation by one frame
// Call this on each tick/frame update
// Returns the current animated value
func (sa *SpringAnimator) Update() float64 {
	sa.value, sa.velocity = sa.spring.Update(sa.value, sa.velocity, sa.target)
	return sa.value
}

// Value returns the current animated value without updating
func (sa *SpringAnimator) Value() float64 {
	return sa.value
}

// IsSettled returns true if the animation has reached the target (within threshold)
func (sa *SpringAnimator) IsSettled() bool {
	const threshold = 0.01
	return abs(sa.target-sa.value) < threshold && abs(sa.velocity) < threshold
}

// Reset sets the value immediately without animation
func (sa *SpringAnimator) Reset(value float64) {
	sa.value = value
	sa.velocity = 0
	sa.target = value
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// ============================================================================
// SMOOTH PROGRESS BAR - Spring-animated progress bar
// ============================================================================

// SmoothProgressBar renders a progress bar with spring-based smooth animation
// animator: SpringAnimator instance (maintains state across calls)
// targetPercent: Target percentage (0.0 to 1.0)
// width: Width of the bar in characters
// Returns: Rendered progress bar string
//
// Example usage:
//
//	type Model struct {
//	    progressAnimator *SpringAnimator
//	    targetProgress   float64
//	}
//
//	func NewModel() Model {
//	    return Model{
//	        progressAnimator: NewSpringAnimator(),
//	        targetProgress:   0.0,
//	    }
//	}
//
//	func (m Model) View() string {
//	    bar := SmoothProgressBar(m.progressAnimator, m.targetProgress, 40)
//	    return fmt.Sprintf("Progress: %s", bar)
//	}
func SmoothProgressBar(animator *SpringAnimator, targetPercent float64, width int) string {
	// Update animator with new target
	animator.SetTarget(targetPercent)
	currentPercent := animator.Update()

	// Calculate filled portion
	filled := int(float64(width) * currentPercent)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	empty := width - filled

	// Build smooth progress bar
	var bar string
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := 0; i < empty; i++ {
		bar += "░"
	}

	return bar
}

// SmoothValue returns a smoothly animated integer value
// Useful for animating counters, percentages, or other numeric displays
//
// Example:
//
//	numberAnimator := NewSpringAnimator()
//	numberAnimator.SetTarget(100.0)
//	for range 60 { // Animate for ~1 second at 60 FPS
//	    value := SmoothValue(numberAnimator)
//	    fmt.Printf("Value: %d\n", value)
//	    time.Sleep(time.Second / 60)
//	}
func SmoothValue(animator *SpringAnimator) int {
	return int(animator.Update() + 0.5) // Round to nearest integer
}

// ============================================================================
// SPRING ANIMATION PRESETS - Pre-configured spring physics for common use cases
// ============================================================================

// NewBouncyAnimator creates a spring animator with bouncy, playful motion
// Perfect for success animations or celebratory moments
// Angular velocity: 7.0, Damping: 0.6 (under-damped = bouncy)
func NewBouncyAnimator() *SpringAnimator {
	return NewSpringAnimatorCustom(7.0, 0.6)
}

// NewSmoothAnimator creates a spring animator with smooth, professional motion
// Perfect for progress bars and professional UI elements
// Angular velocity: 8.0, Damping: 1.0 (critically damped = smooth, no overshoot)
func NewSmoothAnimator() *SpringAnimator {
	return NewSpringAnimatorCustom(8.0, 1.0)
}

// NewSlowAnimator creates a spring animator with slow, deliberate motion
// Perfect for emphasis or important transitions
// Angular velocity: 4.0, Damping: 1.2 (over-damped = slow, no bounce)
func NewSlowAnimator() *SpringAnimator {
	return NewSpringAnimatorCustom(4.0, 1.2)
}

// NewSnappyAnimator creates a spring animator with fast, responsive motion
// Perfect for quick feedback and responsive interactions
// Angular velocity: 12.0, Damping: 1.0 (critically damped but faster)
func NewSnappyAnimator() *SpringAnimator {
	return NewSpringAnimatorCustom(12.0, 1.0)
}
