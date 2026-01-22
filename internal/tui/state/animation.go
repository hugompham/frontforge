package state

// AnimationState holds simple animation counters
// Note: Heavy spring physics animations have been removed per design principles
// Only minimal feedback animations remain (<300ms update cycles)
type AnimationState struct {
	TickCount   int    // Global tick counter for simple animations
	CurrentTask string // Current task description for forging screen
}

// NewAnimationState creates an AnimationState with zero values
func NewAnimationState() AnimationState {
	return AnimationState{
		TickCount:   0,
		CurrentTask: "",
	}
}

// IncrementTick increments the tick counter
func (a *AnimationState) IncrementTick() {
	a.TickCount++
}

// SetTask updates the current task description
func (a *AnimationState) SetTask(task string) {
	a.CurrentTask = task
}

// ResetTick resets the tick counter to zero
func (a *AnimationState) ResetTick() {
	a.TickCount = 0
}
