package tui

import (
	"fmt"
	"frontforge/internal/generators"
	"frontforge/internal/preflight"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// Messages
type runPreflightChecksMsg struct{}
type preflightCompleteMsg struct{}
type generateProjectMsg struct{}
type generationCompleteMsg struct{}
type errorMsg struct{ err error }

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global key handlers (except during critical states)
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			// Prevent accidental quit during preflight checks or generation
			if m.currentState != StatePreflightChecks && m.currentState != StateForging {
				return m, tea.Quit
			}
		}

		// State-specific key handlers
		switch m.currentState {
		case StateTerminalWarn:
			// In warning state, allow user to quit or continue anyway
			switch msg.String() {
			case "enter", " ":
				// User wants to continue despite narrow terminal
				m.currentState = m.previousState
				return m, nil
			case "esc":
				return m, tea.Quit
			}

		case StateWelcome:
			switch msg.String() {
			case "enter", " ":
				// Move from welcome to blueprint phase
				m.currentState = StateBlueprint
				m.anim.TickCount = 0 // Reset tick counter
				// Initialize the form now
				return m, m.form.Init()
			case "esc":
				return m, tea.Quit
			}

		case StateReview:
			switch msg.String() {
			case "f", "F":
				// Show confirmation before starting forge
				m.currentState = StateConfirmForge
				return m, nil
			case "enter":
				// Also allow enter to confirm (same as F)
				m.currentState = StateConfirmForge
				return m, nil
			case "left", "h", "b":
				// Show confirmation before going back
				m.currentState = StateConfirmBack
				return m, nil
			case "esc":
				// Cancel - quit the application
				return m, tea.Quit
			}

		case StateConfirmForge:
			switch msg.String() {
			case "y", "Y", "enter":
				// User confirmed - run pre-flight checks first
				m.applyFormDataToConfig()
				m.currentState = StatePreflightChecks
				m.anim.TickCount = 0 // Reset tick counter
				return m, func() tea.Msg {
					time.Sleep(100 * time.Millisecond) // Give UI time to render
					return runPreflightChecksMsg{}
				}
			case "n", "N", "esc":
				// User cancelled - go back to review
				m.currentState = StateReview
				return m, nil
			}

		case StatePreflightChecks:
			// User can only quit during preflight checks
			switch msg.String() {
			case "esc":
				return m, tea.Quit
			}

		case StateConfirmBack:
			switch msg.String() {
			case "y", "Y", "enter":
				// User confirmed - go back to edit
				// Preserve existing form and its position, just make it editable again
				m.form.State = huh.StateNormal
				m.currentState = StateBlueprint
				return m, nil
			case "n", "N", "esc":
				// User cancelled - stay on review
				m.currentState = StateReview
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.layout.Width = msg.Width
		m.layout.Height = msg.Height
		m.layout.AdaptiveBox = CalculateBoxWidth(msg.Width)

		// Check if terminal is too narrow
		if IsTerminalTooNarrow(msg.Width) && m.currentState != StateTerminalWarn {
			// Save current state and show warning
			m.previousState = m.currentState
			m.currentState = StateTerminalWarn
			return m, nil
		} else if !IsTerminalTooNarrow(msg.Width) && m.currentState == StateTerminalWarn {
			// Terminal is now wide enough, restore previous state
			m.currentState = m.previousState
			return m, nil
		}
		// Note: Huh form width is set at creation and cannot be updated dynamically
		// But we can still adapt other UI elements

	case spinner.TickMsg:
		// Process animation ticks for states that need them
		if m.currentState == StateForging || m.currentState == StateWelcome || m.currentState == StatePreflightChecks {
			// Increment tick counter for simple animations
			m.anim.TickCount++

			// Update spinner and ensure it continues ticking
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		// For all other states, don't tick animations
		return m, nil

	case runPreflightChecksMsg:
		// Run all pre-flight checks
		return m, func() tea.Msg {
			// Add a small delay to make the checks visible
			time.Sleep(500 * time.Millisecond)
			return preflightCompleteMsg{}
		}

	case preflightCompleteMsg:
		// Run the checks
		m.preflightResults = new(preflight.PreflightResults)
		*m.preflightResults = preflight.RunAllChecks(m.config)

		// Check if there were any fatal errors
		if m.preflightResults.FatalError {
			// Show error state with preflight failure details
			m.currentState = StateCracked
			m.err = fmt.Errorf("pre-flight checks failed")
			return m, tea.Quit
		}

		// All checks passed or warnings only - proceed to forging
		m.currentState = StateForging
		m.anim.TickCount = 0 // Reset tick counter for forging animation
		return m, func() tea.Msg {
			time.Sleep(100 * time.Millisecond) // Give UI time to render
			return generateProjectMsg{}
		}

	case generateProjectMsg:
		// Apply form data to config before generation
		m.applyFormDataToConfig()

		// Update task status
		m.anim.CurrentTask = "Generating project structure"

		// Generate the project asynchronously so UI can render
		// Batch the generation command with spinner tick to keep animation running
		return m, tea.Batch(
			func() tea.Msg {
				// Brief animation feedback (500ms max)
				start := time.Now()

				// Generate the project
				err := generators.SetupProject(m.config)

				// Brief minimum display for visual feedback
				elapsed := time.Since(start)
				minDisplayTime := 500 * time.Millisecond
				if elapsed < minDisplayTime {
					time.Sleep(minDisplayTime - elapsed)
				}

				if err != nil {
					return errorMsg{err: err}
				}
				return generationCompleteMsg{}
			},
			m.spinner.Tick, // Keep spinner ticking during generation
		)

	case generationCompleteMsg:
		m.currentState = StateFinished
		return m, tea.Quit

	case errorMsg:
		m.err = msg.err
		m.currentState = StateCracked
		return m, tea.Quit
	}

	// Handle form state
	if m.currentState == StateBlueprint {
		// Update the form with the message
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}

		// Update progress tracker based on key presses (rough approximation)
		if m.progress != nil {
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				switch keyMsg.String() {
				case "enter":
					// Advance progress when user confirms a field
					if m.progress.CurrentStep < m.progress.TotalSteps {
						m.progress.AdvanceStep()
					}
				case "shift+tab":
					// Go back when using shift+tab
					m.progress.GoBack()
				}
			}
		}

		// Check if form is complete using tagged switch
		switch m.form.State {
		case huh.StateCompleted:
			// Move to review screen instead of directly generating
			m.applyFormDataToConfig()
			m.currentState = StateReview
			return m, nil

		case huh.StateAborted:
			// User cancelled
			return m, tea.Quit
		}

		return m, cmd
	}

	// Default: pass to spinner for any unhandled messages (keeps animation running)
	if m.currentState == StateWelcome || m.currentState == StateForging || m.currentState == StatePreflightChecks {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}
