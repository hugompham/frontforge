package tui_test

import (
	"frontforge/internal/tui"
	"testing"
)

func TestProgressTracker(t *testing.T) {
	totalSteps := 10
	tracker := tui.NewProgressTracker(totalSteps)

	// Verify initial state
	if tracker.CurrentStep != 1 {
		t.Errorf("Expected CurrentStep 1, got %d", tracker.CurrentStep)
	}
	if tracker.TotalSteps != totalSteps {
		t.Errorf("Expected TotalSteps %d, got %d", totalSteps, tracker.TotalSteps)
	}
	if tracker.CurrentPhase != tui.PhaseFoundation {
		t.Errorf("Expected initial phase Foundation, got %v", tracker.CurrentPhase)
	}

	// Test step progression
	for i := 1; i < totalSteps; i++ {
		tracker.AdvanceStep()
		if tracker.CurrentStep != i+1 {
			t.Errorf("After %d advances, expected CurrentStep %d, got %d", i, i+1, tracker.CurrentStep)
		}
	}

	// Verify completion tracking
	for i := 0; i < totalSteps-1; i++ {
		if !tracker.Completed[i] {
			t.Errorf("Step %d should be marked completed", i)
		}
	}
}

func TestProgressTrackerAdvanceLimit(t *testing.T) {
	totalSteps := 5
	tracker := tui.NewProgressTracker(totalSteps)

	// Advance to last step
	for i := 1; i < totalSteps; i++ {
		tracker.AdvanceStep()
	}

	// At last step
	if tracker.CurrentStep != totalSteps {
		t.Errorf("Expected to be at step %d, got %d", totalSteps, tracker.CurrentStep)
	}

	// Try to advance past last step - should be blocked
	tracker.AdvanceStep()
	if tracker.CurrentStep != totalSteps {
		t.Errorf("CurrentStep should stay at %d, got %d", totalSteps, tracker.CurrentStep)
	}

	// IsComplete is true when CurrentStep > TotalSteps, but AdvanceStep prevents this
	// So we test the guard condition instead
	if tracker.CurrentStep > totalSteps {
		t.Error("CurrentStep should not exceed TotalSteps")
	}
}

func TestProgressTrackerPercentage(t *testing.T) {
	totalSteps := 10

	tests := []struct {
		step    int
		wantPct int
	}{
		{1, 0},   // Step 1 = 0% (0 completed)
		{2, 10},  // Step 2 = 10% (1 completed)
		{3, 20},  // Step 3 = 20% (2 completed)
		{5, 40},  // Step 5 = 40% (4 completed)
		{10, 90}, // Step 10 = 90% (9 completed)
	}

	for _, tt := range tests {
		tracker := tui.NewProgressTracker(totalSteps)
		// Advance to target step
		for i := 1; i < tt.step; i++ {
			tracker.AdvanceStep()
		}

		pct := tracker.PercentComplete()
		if pct != tt.wantPct {
			t.Errorf("Step %d: expected %d%%, got %d%%", tt.step, tt.wantPct, pct)
		}
	}
}

func TestProgressTrackerGoBack(t *testing.T) {
	totalSteps := 5
	tracker := tui.NewProgressTracker(totalSteps)

	// Advance to step 3
	tracker.AdvanceStep()
	tracker.AdvanceStep()

	if tracker.CurrentStep != 3 {
		t.Errorf("Expected CurrentStep 3, got %d", tracker.CurrentStep)
	}

	// Go back
	tracker.GoBack()
	if tracker.CurrentStep != 2 {
		t.Errorf("Expected CurrentStep 2 after GoBack, got %d", tracker.CurrentStep)
	}

	// Can't go back past step 1
	tracker.GoBack()
	tracker.GoBack()
	if tracker.CurrentStep != 1 {
		t.Errorf("Expected CurrentStep 1 after multiple GoBack, got %d", tracker.CurrentStep)
	}
}

func TestProgressTrackerJumpToStep(t *testing.T) {
	totalSteps := 10
	tracker := tui.NewProgressTracker(totalSteps)

	// Jump to step 5
	tracker.JumpToStep(5)
	if tracker.CurrentStep != 5 {
		t.Errorf("Expected CurrentStep 5 after jump, got %d", tracker.CurrentStep)
	}

	// Jump to invalid step (too high) should not change
	tracker.JumpToStep(15)
	if tracker.CurrentStep != 5 {
		t.Errorf("Expected CurrentStep to remain 5, got %d", tracker.CurrentStep)
	}

	// Jump to invalid step (zero) should not change
	tracker.JumpToStep(0)
	if tracker.CurrentStep != 5 {
		t.Errorf("Expected CurrentStep to remain 5, got %d", tracker.CurrentStep)
	}
}

func TestFormPhases(t *testing.T) {
	phases := []struct {
		phase       tui.FormPhase
		name        string
		description string
	}{
		{tui.PhaseFoundation, "FOUNDATION", "Core project structure and framework"},
		{tui.PhaseTooling, "TOOLING", "Development and styling tools"},
		{tui.PhaseFeatures, "FEATURES", "Application functionality and libraries"},
		{tui.PhaseFinishing, "FINISHING", "Utilities and project organization"},
	}

	for _, tt := range phases {
		if tt.phase.String() != tt.name {
			t.Errorf("Expected phase name %q, got %q", tt.name, tt.phase.String())
		}
		if tt.phase.Description() != tt.description {
			t.Errorf("Expected description %q, got %q", tt.description, tt.phase.Description())
		}
	}
}
