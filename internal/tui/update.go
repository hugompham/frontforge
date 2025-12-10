package tui

import (
	"frontforge/internal/generators"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

// Messages
type generateProjectMsg struct{}
type generationCompleteMsg struct{}
type errorMsg struct{ err error }

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global key handlers
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		// State-specific key handlers
		switch m.state {
		case StateTerminalWarn:
			// In warning state, allow user to quit or continue anyway
			switch msg.String() {
			case "enter", " ":
				// User wants to continue despite narrow terminal
				m.state = m.previousState
				return m, nil
			case "esc":
				return m, tea.Quit
			}

		case StateWelcome:
			switch msg.String() {
			case "enter", " ":
				// Move from welcome to blueprint phase
				m.state = StateBlueprint
				m.tickCount = 0 // Reset tick counter
				// Initialize the form now
				return m, m.form.Init()
			case "esc":
				return m, tea.Quit
			}

		case StateReview:
			switch msg.String() {
			case "f", "F":
				// Show confirmation before starting forge
				m.state = StateConfirmForge
				return m, nil
			case "enter":
				// Also allow enter to confirm (same as F)
				m.state = StateConfirmForge
				return m, nil
			case "left", "h", "b":
				// Show confirmation before going back
				m.state = StateConfirmBack
				return m, nil
			case "esc":
				// Cancel - quit the application
				return m, tea.Quit
			}

		case StateConfirmForge:
			switch msg.String() {
			case "y", "Y", "enter":
				// User confirmed - start forging
				m.state = StateForging
				m.tickCount = 0 // Reset tick counter for forging animation
				m.applyFormDataToConfig()
				// Add small delay so UI can render the forging screen first
				return m, func() tea.Msg {
					time.Sleep(100 * time.Millisecond) // Give UI time to render
					return generateProjectMsg{}
				}
			case "n", "N", "esc":
				// User cancelled - go back to review
				m.state = StateReview
				return m, nil
			}

		case StateConfirmBack:
			switch msg.String() {
			case "y", "Y", "enter":
				// User confirmed - go back to edit
				m.form = m.createForm()
				m.form.State = huh.StateNormal
				m.state = StateBlueprint
				return m, m.form.Init()
			case "n", "N", "esc":
				// User cancelled - stay on review
				m.state = StateReview
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.adaptiveBox = CalculateBoxWidth(msg.Width)

		// Check if terminal is too narrow
		if IsTerminalTooNarrow(msg.Width) && m.state != StateTerminalWarn {
			// Save current state and show warning
			m.previousState = m.state
			m.state = StateTerminalWarn
			return m, nil
		} else if !IsTerminalTooNarrow(msg.Width) && m.state == StateTerminalWarn {
			// Terminal is now wide enough, restore previous state
			m.state = m.previousState
			return m, nil
		}
		// Note: Huh form width is set at creation and cannot be updated dynamically
		// But we can still adapt other UI elements

	case spinner.TickMsg:
		// Process animation ticks for states that need them
		if m.state == StateForging || m.state == StateWelcome {
			// Increment tick counter for animations
			m.tickCount++

			// Update flame animation with spring physics (welcome state)
			if m.state == StateWelcome {
				// Toggle flame target every 3600 ticks (~60 seconds / 1 minute at 60fps) - extremely slow breathing
				if m.tickCount%3600 == 0 {
					m.flameTargetToggle = !m.flameTargetToggle
				}

				// Set target based on toggle (oscillate between low and high flames)
				target := 0.35 // Low flames
				if m.flameTargetToggle {
					target = 0.75 // High flames (reduced range for subtlety)
				}

				// Update spring physics - this is the key Harmonica method
				m.flameIntensity, m.flameVelocity = m.flameSpring.Update(
					m.flameIntensity,
					m.flameVelocity,
					target,
				)
			}

			// Update forging hammer animation with spring physics (forging state)
			if m.state == StateForging {
				// Toggle hammer target every 3600 ticks (~60 seconds at 60fps) - exactly match furnace speed
				if m.tickCount%3600 == 0 {
					m.forgingFlameTargetToggle = !m.forgingFlameTargetToggle
				}

				// Set target based on toggle (hammer position)
				target := 0.1 // Hammer raised high
				if m.forgingFlameTargetToggle {
					target = 0.95 // Hammer striking anvil
				}

				// Update spring physics for hammer animation every tick for smooth motion
				m.forgingFlameIntensity, m.forgingFlameVelocity = m.forgingFlameSpring.Update(
					m.forgingFlameIntensity,
					m.forgingFlameVelocity,
					target,
				)
			}

			// Update spinner and ensure it continues ticking
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		// For all other states, don't tick animations
		return m, nil

	case generateProjectMsg:
		// Apply form data to config before generation
		m.applyFormDataToConfig()

		// Update task status
		m.currentTask = "Generating project structure"

		// Generate the project asynchronously so UI can render
		// Batch the generation command with spinner tick to keep animation running
		return m, tea.Batch(
			func() tea.Msg {
				// Ensure animation is visible (5 seconds for forging animation)
				start := time.Now()

				// Generate the project
				err := generators.SetupProject(m.config)

				// Ensure minimum display time for animation
				elapsed := time.Since(start)
				minDisplayTime := 5 * time.Second
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
		m.state = StateFinished
		return m, tea.Quit

	case errorMsg:
		m.err = msg.err
		m.state = StateCracked
		return m, tea.Quit
	}

	// Handle form state
	if m.state == StateBlueprint {
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
			m.state = StateReview
			return m, nil

		case huh.StateAborted:
			// User cancelled
			return m, tea.Quit
		}

		return m, cmd
	}

	// Default: pass to spinner for any unhandled messages (keeps animation running)
	if m.state == StateWelcome || m.state == StateForging {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}
