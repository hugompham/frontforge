package tui_test

import (
	"frontforge/internal/models"
	"frontforge/internal/tui"
	"frontforge/internal/tui/state"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	m := tui.NewModel()

	// Verify initial state
	if m.GetCurrentState() != tui.StateWelcome {
		t.Errorf("Expected initial state to be StateWelcome, got %v", m.GetCurrentState())
	}

	// Verify form state has defaults
	formState := m.GetFormState()
	if formState.Language != "TypeScript" {
		t.Errorf("Expected default language TypeScript, got %s", formState.Language)
	}
	if formState.Framework != "React" {
		t.Errorf("Expected default framework React, got %s", formState.Framework)
	}
	if formState.PackageManager != "npm" {
		t.Errorf("Expected default package manager npm, got %s", formState.PackageManager)
	}
}

func TestStateTransitions(t *testing.T) {
	tests := []struct {
		name        string
		fromState   tui.State
		msg         tea.Msg
		expectState tui.State
	}{
		{
			name:        "Welcome to Blueprint on Enter",
			fromState:   tui.StateWelcome,
			msg:         tea.KeyMsg{Type: tea.KeyEnter},
			expectState: tui.StateBlueprint,
		},
		{
			name:        "Welcome to Quit on q",
			fromState:   tui.StateWelcome,
			msg:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			expectState: tui.StateWelcome, // Will actually quit, but state doesn't change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tui.NewModel()
			m.SetCurrentState(tt.fromState)

			updatedModel, _ := m.Update(tt.msg)
			m = updatedModel.(tui.Model)

			if m.GetCurrentState() != tt.expectState {
				t.Errorf("Expected state %v, got %v", tt.expectState, m.GetCurrentState())
			}
		})
	}
}

func TestProjectNameValidation(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantValid bool
		wantError string
	}{
		{
			name:      "Valid lowercase",
			input:     "my-app",
			wantValid: true,
		},
		{
			name:      "Valid with numbers",
			input:     "app123",
			wantValid: true,
		},
		{
			name:      "Valid with underscore",
			input:     "my_app",
			wantValid: true,
		},
		{
			name:      "Valid mixed case",
			input:     "MyApp",
			wantValid: true,
		},
		{
			name:      "Invalid with spaces",
			input:     "my app",
			wantValid: false,
			wantError: "alphanumeric",
		},
		{
			name:      "Invalid with special chars",
			input:     "my@app",
			wantValid: false,
			wantError: "alphanumeric",
		},
		{
			name:      "Invalid empty",
			input:     "",
			wantValid: false,
			wantError: "required",
		},
		{
			name:      "Invalid with slash",
			input:     "my/app",
			wantValid: false,
			wantError: "alphanumeric",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tui.ValidateProjectName(tt.input)
			if tt.wantValid {
				if err != nil {
					t.Errorf("Expected valid, got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.wantError)
				}
			}
		})
	}
}

func TestFormStateDefaults(t *testing.T) {
	formState := state.NewFormState()

	// Verify all defaults are set
	defaults := map[string]string{
		"SetupMode":       "custom",
		"ProjectName":     "my-app",
		"Language":        "TypeScript",
		"Framework":       "React",
		"PackageManager":  "npm",
		"Styling":         "Tailwind CSS",
		"UILibrary":       "Shadcn/ui",
		"Routing":         "React Router",
		"Testing":         "Vitest",
		"StateManagement": "Zustand",
		"FormManagement":  "React Hook Form",
		"DataFetching":    "TanStack Query",
		"Animation":       "Framer Motion",
		"Icons":           "Heroicons",
		"DataViz":         "None",
		"Utilities":       "date-fns",
		"I18n":            "None",
		"Structure":       "Feature-based",
	}

	// Check each field
	if formState.SetupMode != defaults["SetupMode"] {
		t.Errorf("SetupMode: got %q, want %q", formState.SetupMode, defaults["SetupMode"])
	}
	if formState.ProjectName != defaults["ProjectName"] {
		t.Errorf("ProjectName: got %q, want %q", formState.ProjectName, defaults["ProjectName"])
	}
	if formState.Language != defaults["Language"] {
		t.Errorf("Language: got %q, want %q", formState.Language, defaults["Language"])
	}
	if formState.Framework != defaults["Framework"] {
		t.Errorf("Framework: got %q, want %q", formState.Framework, defaults["Framework"])
	}
	if formState.PackageManager != defaults["PackageManager"] {
		t.Errorf("PackageManager: got %q, want %q", formState.PackageManager, defaults["PackageManager"])
	}
	if formState.Styling != defaults["Styling"] {
		t.Errorf("Styling: got %q, want %q", formState.Styling, defaults["Styling"])
	}
}

func TestConfigMapping(t *testing.T) {
	tests := []struct {
		name      string
		formState state.FormState
		want      models.Config
	}{
		{
			name: "React TypeScript with Tailwind",
			formState: state.FormState{
				ProjectName:     "test-app",
				Language:        "TypeScript",
				Framework:       "React",
				PackageManager:  "npm",
				Styling:         "Tailwind CSS",
				Routing:         "React Router",
				Testing:         "Vitest",
				StateManagement: "Zustand",
				DataFetching:    "TanStack Query",
				Structure:       "Feature-based",
			},
			want: models.Config{
				ProjectName:     "test-app",
				Language:        "TypeScript",
				Framework:       "React",
				PackageManager:  "npm",
				Styling:         "Tailwind CSS",
				Routing:         "React Router",
				Testing:         "Vitest",
				StateManagement: "Zustand",
				DataFetching:    "TanStack Query",
				Structure:       "Feature-based",
			},
		},
		{
			name: "Vue JavaScript with CSS Modules",
			formState: state.FormState{
				ProjectName:    "vue-app",
				Language:       "JavaScript",
				Framework:      "Vue",
				PackageManager: "pnpm",
				Styling:        "CSS Modules",
				Routing:        "None",
				Testing:        "None",
				Structure:      "Layer-based",
			},
			want: models.Config{
				ProjectName:    "vue-app",
				Language:       "JavaScript",
				Framework:      "Vue",
				PackageManager: "pnpm",
				Styling:        "CSS Modules",
				Routing:        "None",
				Testing:        "None",
				Structure:      "Layer-based",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tui.NewModel()
			m.SetFormState(tt.formState)
			m.ApplyFormDataToConfig()
			config := m.GetConfig()

			if config.ProjectName != tt.want.ProjectName {
				t.Errorf("ProjectName: got %q, want %q", config.ProjectName, tt.want.ProjectName)
			}
			if config.Language != tt.want.Language {
				t.Errorf("Language: got %q, want %q", config.Language, tt.want.Language)
			}
			if config.Framework != tt.want.Framework {
				t.Errorf("Framework: got %q, want %q", config.Framework, tt.want.Framework)
			}
			if config.PackageManager != tt.want.PackageManager {
				t.Errorf("PackageManager: got %q, want %q", config.PackageManager, tt.want.PackageManager)
			}
			if config.Styling != tt.want.Styling {
				t.Errorf("Styling: got %q, want %q", config.Styling, tt.want.Styling)
			}
		})
	}
}

func TestStateAliases(t *testing.T) {
	// Verify legacy state aliases work correctly
	if tui.StateForm != tui.StateBlueprint {
		t.Error("StateForm should alias StateBlueprint")
	}
	if tui.StateGenerating != tui.StateForging {
		t.Error("StateGenerating should alias StateForging")
	}
	if tui.StateSuccess != tui.StateFinished {
		t.Error("StateSuccess should alias StateFinished")
	}
	if tui.StateError != tui.StateCracked {
		t.Error("StateError should alias StateCracked")
	}
}

func TestLayoutState(t *testing.T) {
	layoutState := state.NewLayoutState()

	// Verify reasonable defaults
	if layoutState.Width != 80 {
		t.Errorf("Expected default width 80, got %d", layoutState.Width)
	}
	if layoutState.Height != 24 {
		t.Errorf("Expected default height 24, got %d", layoutState.Height)
	}
	if layoutState.AdaptiveBox != 76 {
		t.Errorf("Expected default adaptive box 76, got %d", layoutState.AdaptiveBox)
	}
}

func TestAnimationState(t *testing.T) {
	animState := state.NewAnimationState()

	// Verify animation state initialized
	if animState.TickCount != 0 {
		t.Errorf("Expected TickCount 0, got %d", animState.TickCount)
	}
	if animState.CurrentTask != "" {
		t.Errorf("Expected empty CurrentTask, got %q", animState.CurrentTask)
	}

	// Test incrementing
	animState.IncrementTick()
	if animState.TickCount != 1 {
		t.Errorf("Expected TickCount 1 after increment, got %d", animState.TickCount)
	}

	// Test setting task
	animState.SetTask("Test task")
	if animState.CurrentTask != "Test task" {
		t.Errorf("Expected CurrentTask 'Test task', got %q", animState.CurrentTask)
	}

	// Test reset
	animState.ResetTick()
	if animState.TickCount != 0 {
		t.Errorf("Expected TickCount 0 after reset, got %d", animState.TickCount)
	}
}
