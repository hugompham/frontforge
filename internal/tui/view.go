package tui

import (
	"fmt"
	"frontforge/internal/models"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current state
// All views use adaptive width calculated from m.adaptiveBox for responsiveness
// Standard pattern: contentWidth := m.adaptiveBox - 4 (accounting for padding)
func (m Model) View() string {
	switch m.state {
	case StateTerminalWarn:
		return m.viewTerminalWarning()
	case StateWelcome:
		return m.viewWelcome()
	case StateBlueprint: // StateForm is alias to StateBlueprint
		return m.viewForm()
	case StateReview:
		return m.viewReview()
	case StateConfirmForge:
		return m.viewConfirmForge()
	case StateConfirmBack:
		return m.viewConfirmBack()
	case StateForging: // StateGenerating is alias to StateForging
		return m.viewGenerating()
	case StateFinished: // StateSuccess is alias to StateFinished
		return m.viewSuccess()
	case StateCracked: // StateError is alias to StateCracked
		return m.viewError()
	default:
		return ""
	}
}

// viewTerminalWarning renders a warning when terminal is too narrow
func (m Model) viewTerminalWarning() string {
	var b strings.Builder

	// Warning header with icon and color - accessible to color-blind users
	header := RenderStatusHeader("warning", "TERMINAL TOO NARROW")
	b.WriteString(header + "\n\n")

	// Warning message
	warningMsg := lipgloss.NewStyle().
		Foreground(colorDraftPencil).
		Width(minBoxWidth).
		Render("Your terminal window is too narrow for optimal display.")
	b.WriteString(warningMsg + "\n\n")

	// Current dimensions
	dimensionsStyle := lipgloss.NewStyle().
		Foreground(colorAnvilGray)
	b.WriteString(dimensionsStyle.Render(fmt.Sprintf("Current width: %d characters", m.width)) + "\n")
	b.WriteString(dimensionsStyle.Render("Recommended: 80 characters or more") + "\n\n")

	// Divider
	dividerWidth := minBoxWidth
	divider := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render(strings.Repeat("─", dividerWidth))
	b.WriteString(divider + "\n\n")

	// Instructions
	instructionsStyle := lipgloss.NewStyle().
		Foreground(colorSilverSheen).
		Bold(true)
	b.WriteString(instructionsStyle.Render("WHAT TO DO:") + "\n\n")

	helpStyle := lipgloss.NewStyle().
		Foreground(colorDraftPencil)
	b.WriteString(helpStyle.Render("1. Resize your terminal window to be wider") + "\n")
	b.WriteString(helpStyle.Render("2. Or press ENTER to continue anyway (may have display issues)") + "\n")
	b.WriteString(helpStyle.Render("3. Or press ESC to exit") + "\n\n")

	b.WriteString(divider + "\n\n")

	// Exit hint
	exitStyle := lipgloss.NewStyle().
		Foreground(colorAshGray).
		Align(lipgloss.Center).
		Width(minBoxWidth)
	b.WriteString(exitStyle.Render("The application will resume automatically when you resize") + "\n")

	return lipgloss.NewStyle().
		Padding(2, 0).
		Render(b.String())
}

// viewWelcome renders the welcome/introduction screen
func (m Model) viewWelcome() string {
	var b strings.Builder

	// Use adaptive width
	contentWidth := m.adaptiveBox - 4

	// Title header
	title := GetForgeTitle()
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorFurnaceOrange).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(titleStyle.Render(title) + "\n")

	// Subtitle
	subtitle := "Craft Your Frontend Project in the Forge"
	subtitleStyle := lipgloss.NewStyle().
		Foreground(colorSilverSheen).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(subtitleStyle.Render(subtitle) + "\n\n")

	// Forge furnace animation with breathing flames (using Harmonica spring physics)
	animation := GetWelcomeAnimationFrame(m.flameIntensity)
	animStyle := lipgloss.NewStyle().
		Foreground(colorMoltenGold).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(animStyle.Render(animation) + "\n")

	// Adaptive divider
	divider := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render(strings.Repeat("─", contentWidth))
	b.WriteString(divider + "\n\n")

	// Introduction text
	intro := "Configure your project through four phases: Foundation, Tooling, Features, and Finishing."
	introStyle := lipgloss.NewStyle().
		Foreground(colorDraftPencil).
		Width(contentWidth)
	b.WriteString(introStyle.Render(intro) + "\n\n")

	// Call to action
	ctaStyle := lipgloss.NewStyle().
		Foreground(colorMoltenGold).
		Bold(true).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(ctaStyle.Render("Press ENTER to begin") + "\n\n")

	// Exit hint
	exitStyle := lipgloss.NewStyle().
		Foreground(colorAshGray).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(exitStyle.Render("Press ESC to exit") + "\n")

	// Return without border/card wrapper
	return lipgloss.NewStyle().
		Padding(2, 0).
		Render(b.String())
}

func (m Model) viewForm() string {
	// Just show the form - Huh handles everything
	return m.form.View()
}

// viewReview renders the configuration review screen before forging
func (m Model) viewReview() string {
	var b strings.Builder
	contentWidth := m.adaptiveBox - 4

	// Enhanced header
	headerText := "BLUEPRINT COMPLETE"
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorFurnaceOrange).
		Background(colorIronDark).
		Padding(0, 2).
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(headerText)
	b.WriteString(header + "\n\n")

	// Subtitle explaining what this screen is
	subtitle := lipgloss.NewStyle().
		Foreground(colorDraftPencil).
		Align(lipgloss.Center).
		Width(contentWidth).
		Render("Review your configuration before forging begins")
	b.WriteString(subtitle + "\n\n")

	// Adaptive divider
	divider := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render(strings.Repeat("─", contentWidth))
	b.WriteString(divider + "\n\n")

	// Enhanced configuration display with sections
	b.WriteString(m.getEnhancedConfig() + "\n")

	b.WriteString(divider + "\n\n")

	// Action buttons section with clear visual hierarchy
	actionsTitle := lipgloss.NewStyle().
		Foreground(colorSilverSheen).
		Bold(true).
		Render("ACTIONS")
	b.WriteString(actionsTitle + "\n\n")

	// Primary action - Start forging
	primaryAction := lipgloss.NewStyle().
		Foreground(colorCharcoal).
		Background(colorFurnaceOrange).
		Bold(true).
		Padding(0, 2).
		Render("START FORGING")

	// Secondary actions
	backAction := lipgloss.NewStyle().
		Foreground(colorSteelBlue).
		Render("[←] Back to edit")

	cancelAction := lipgloss.NewStyle().
		Foreground(colorAshGray).
		Render("[Esc] Cancel")

	// Center the primary action
	primaryCentered := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(primaryAction)
	b.WriteString(primaryCentered + "\n\n")

	// Show secondary actions side by side
	secondaryActions := backAction + "  " + lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render("•") + "  " + cancelAction
	secondaryCentered := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(secondaryActions)
	b.WriteString(secondaryCentered + "\n")

	// Return without border
	return lipgloss.NewStyle().
		Padding(1, 0).
		Render(b.String())
}

func (m Model) viewGenerating() string {
	var b strings.Builder
	contentWidth := m.adaptiveBox - 4

	// Header with forging theme
	header := forgeHeaderStyle.Render("F O R G I N G")
	b.WriteString(header + "\n\n")

	subtitle := helpTextStyle.Render("Your project is being shaped in the forge...")
	b.WriteString(subtitle + "\n\n")

	// Animated forging visualization - hammer striking anvil with spring physics
	animation := GetForgingAnimationFrame(m.forgingFlameIntensity)
	animStyle := lipgloss.NewStyle().
		Foreground(colorMoltenGold).
		Bold(true).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(animStyle.Render(animation) + "\n\n")

	// Current task with animated dots
	if m.currentTask != "" {
		task := forgeMutedStyle.Render("Current task: ") +
			forgeInfoStyle.Render(m.currentTask) +
			AnimatedDots(m.tickCount, 3)
		b.WriteString(task + "\n\n")
	} else {
		task := forgeMutedStyle.Render("Setting up your project structure") +
			AnimatedDots(m.tickCount, 3)
		b.WriteString(task + "\n\n")
	}

	// Adaptive divider
	divider := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render(strings.Repeat("─", contentWidth))
	b.WriteString(divider + "\n\n")

	// Show configuration being forged
	b.WriteString(phaseTitleStyle.Render("  YOUR BLUEPRINT  ") + "\n\n")
	config := m.getStyledConfig()
	b.WriteString(config + "\n")

	// Return without border
	return lipgloss.NewStyle().
		Padding(1, 0).
		Render(b.String())
}

func (m Model) viewSuccess() string {
	var b strings.Builder
	contentWidth := m.adaptiveBox - 4

	// Success header with icon and color - accessible to color-blind users
	header := RenderStatusHeader("success", "PROJECT CREATED SUCCESSFULLY")
	b.WriteString(header + "\n\n")

	// Celebration subtitle
	subtitle := lipgloss.NewStyle().
		Foreground(colorTemperedGreen).
		Bold(true).
		Render("Your project is ready to use!")
	b.WriteString(subtitle + "\n\n")

	// Adaptive divider
	divider := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render(strings.Repeat("─", contentWidth))
	b.WriteString(divider + "\n\n")

	// Project summary section
	b.WriteString(phaseTitleStyle.Render("  PROJECT SUMMARY  ") + "\n\n")

	// Project location
	b.WriteString(forgeLabelStyle.Render("Location:"))
	b.WriteString(inputValueStyle.Render(m.config.ProjectPath) + "\n\n")

	// Configuration summary
	summary := m.getStyledConfig()
	b.WriteString(summary + "\n")

	b.WriteString(divider + "\n\n")

	// Next steps section with clear instructions
	b.WriteString(sectionHeaderStyle.Render("NEXT STEPS") + "\n\n")

	// Build command strings
	// Check if project is in current directory or a subfolder
	cwd, _ := filepath.Abs(".")
	isCurrentDir := m.config.ProjectPath == cwd

	cdCmd := fmt.Sprintf("cd %s", m.config.ProjectName)
	installCmd := fmt.Sprintf("%s install", m.config.PackageManager)
	devCmd := fmt.Sprintf("%s ", m.config.PackageManager)
	if m.config.PackageManager == "npm" {
		devCmd += "run "
	}
	devCmd += "dev"

	// Copyable commands section (plain text for easy copy)
	copyHintStyle := lipgloss.NewStyle().
		Foreground(colorDraftPencil).
		Italic(false)
	b.WriteString(copyHintStyle.Render("Copy these commands:") + "\n\n")

	plainCmdStyle := lipgloss.NewStyle().
		Foreground(colorSilverSheen)

	// Only show cd command if not in current directory
	if !isCurrentDir {
		b.WriteString(plainCmdStyle.Render("  " + cdCmd) + "\n")
	}
	b.WriteString(plainCmdStyle.Render("  " + installCmd) + "\n")
	b.WriteString(plainCmdStyle.Render("  " + devCmd) + "\n\n")

	// Styled command boxes for visual appeal
	b.WriteString(copyHintStyle.Render("Or follow these steps:") + "\n\n")

	stepNumStyle := lipgloss.NewStyle().
		Foreground(colorMoltenGold).
		Bold(true)
	stepNum := 1

	// Step 1: Navigate to project (only if not in current directory)
	if !isCurrentDir {
		b.WriteString(stepNumStyle.Render(fmt.Sprintf("%d. ", stepNum)) + forgeLabelStyle.Render("Navigate to your project") + "\n")
		cmdBox := lipgloss.NewStyle().
			Foreground(colorMoltenGold).
			Background(lipgloss.Color("#2D2D2D")).
			Padding(0, 1).
			Render(cdCmd)
		b.WriteString("   " + cmdBox + "\n\n")
		stepNum++
	}

	// Step 2 (or 1 if current dir): Install dependencies
	b.WriteString(stepNumStyle.Render(fmt.Sprintf("%d. ", stepNum)) + forgeLabelStyle.Render("Install dependencies") + "\n")
	stepNum++
	installBox := lipgloss.NewStyle().
		Foreground(colorMoltenGold).
		Background(lipgloss.Color("#2D2D2D")).
		Padding(0, 1).
		Render(installCmd)
	b.WriteString("   " + installBox + "\n\n")

	// Step 3 (or 2 if current dir): Start dev server
	b.WriteString(stepNumStyle.Render(fmt.Sprintf("%d. ", stepNum)) + forgeLabelStyle.Render("Start development server") + "\n")
	devBox := lipgloss.NewStyle().
		Foreground(colorMoltenGold).
		Background(lipgloss.Color("#2D2D2D")).
		Padding(0, 1).
		Render(devCmd)
	b.WriteString("   " + devBox + "\n\n")

	b.WriteString(divider + "\n\n")

	// Closing message with helpful tip
	closingStyle := lipgloss.NewStyle().
		Foreground(colorSilverSheen).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(closingStyle.Render("Happy coding!") + "\n")

	// Return without border
	return lipgloss.NewStyle().
		Padding(1, 0).
		Render(b.String())
}

func (m Model) viewError() string {
	var b strings.Builder
	contentWidth := m.adaptiveBox - 4

	errMsg := "An unknown error occurred"
	if m.err != nil {
		errMsg = m.err.Error()
	}

	// Error header with icon and color - accessible to color-blind users
	header := RenderStatusHeader("error", "PROJECT CREATION FAILED")
	b.WriteString(header + "\n\n")

	subtitle := helpTextStyle.Render("Something went wrong during project generation.")
	b.WriteString(subtitle + "\n\n")

	// Adaptive divider
	divider := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render(strings.Repeat("─", contentWidth))
	b.WriteString(divider + "\n\n")

	// Error details section
	b.WriteString(sectionHeaderStyle.Render("ERROR DETAILS") + "\n\n")

	// Error message box with better visibility
	errorBoxWidth := m.adaptiveBox - 8
	if errorBoxWidth < 40 {
		errorBoxWidth = 40
	}
	errorBox := lipgloss.NewStyle().
		Foreground(colorCrackedRed).
		Background(lipgloss.Color("#2D2D2D")).
		Padding(1, 2).
		Width(errorBoxWidth).
		Render(errMsg)
	b.WriteString(errorBox + "\n\n")

	b.WriteString(divider + "\n\n")

	// Troubleshooting section
	b.WriteString(sectionHeaderStyle.Render("TROUBLESHOOTING") + "\n\n")
	b.WriteString(forgeLabelStyle.Render("Common solutions:") + "\n")
	b.WriteString(forgeMutedStyle.Render("  - Ensure you have a stable internet connection") + "\n")
	b.WriteString(forgeMutedStyle.Render("  - Verify required tools are installed (Node.js, npm/yarn/pnpm/bun)") + "\n")
	b.WriteString(forgeMutedStyle.Render("  - Check that the target directory doesn't already exist") + "\n")
	b.WriteString(forgeMutedStyle.Render("  - Try running with administrator/sudo privileges if needed") + "\n\n")

	b.WriteString(divider + "\n\n")

	// Exit instruction
	exitStyle := lipgloss.NewStyle().
		Foreground(colorAshGray).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(exitStyle.Render("Press ESC or Ctrl+C to exit") + "\n")

	// Return without border
	return lipgloss.NewStyle().
		Padding(1, 0).
		Render(b.String())
}

// getStyledConfig returns a formatted config summary (NO letter icons)
func (m Model) getStyledConfig() string {
	var b strings.Builder

	// Framework (no icon)
	b.WriteString(forgeLabelStyle.Render("Framework:"))
	b.WriteString(inputValueStyle.Render(m.config.Framework) + "\n")

	// Language
	b.WriteString(forgeLabelStyle.Render("Language:"))
	b.WriteString(inputValueStyle.Render(m.config.Language) + "\n")

	// Package Manager (no icon)
	b.WriteString(forgeLabelStyle.Render("Package Manager:"))
	b.WriteString(inputValueStyle.Render(m.config.PackageManager) + "\n")

	// Styling (no icon)
	b.WriteString(forgeLabelStyle.Render("Styling:"))
	b.WriteString(inputValueStyle.Render(m.config.Styling) + "\n")

	// UI Library (if any, no icon)
	if m.config.UILibrary != "" && m.config.UILibrary != models.UILibraryNone {
		b.WriteString(forgeLabelStyle.Render("UI Library:"))
		b.WriteString(inputValueStyle.Render(m.config.UILibrary) + "\n")
	}

	// Testing (if any, no icon)
	if m.config.Testing != "" && m.config.Testing != models.TestingNone {
		b.WriteString(forgeLabelStyle.Render("Testing:"))
		b.WriteString(inputValueStyle.Render(m.config.Testing) + "\n")
	}

	// State Management (if any, no icon)
	if m.config.StateManagement != "" && m.config.StateManagement != models.StateNone {
		b.WriteString(forgeLabelStyle.Render("State Management:"))
		b.WriteString(inputValueStyle.Render(m.config.StateManagement) + "\n")
	}

	// Data Fetching (if any, no icon)
	if m.config.DataFetching != "" && m.config.DataFetching != models.DataNone {
		b.WriteString(forgeLabelStyle.Render("Data Fetching:"))
		b.WriteString(inputValueStyle.Render(m.config.DataFetching) + "\n")
	}

	return b.String()
}

// getCompactConfig returns a compact configuration summary for the review screen
func (m Model) getCompactConfig() string {
	var b strings.Builder

	// Compact label style (shorter width)
	compactLabelStyle := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Width(14)

	// Project name
	b.WriteString(compactLabelStyle.Render("Project:"))
	b.WriteString(inputValueStyle.Render(m.config.ProjectName) + "\n")

	// Framework
	b.WriteString(compactLabelStyle.Render("Framework:"))
	b.WriteString(inputValueStyle.Render(m.config.Framework) + "\n")

	// Language
	b.WriteString(compactLabelStyle.Render("Language:"))
	b.WriteString(inputValueStyle.Render(m.config.Language) + "\n")

	// Package Manager
	b.WriteString(compactLabelStyle.Render("Pkg Manager:"))
	b.WriteString(inputValueStyle.Render(m.config.PackageManager) + "\n")

	// Styling
	b.WriteString(compactLabelStyle.Render("Styling:"))
	b.WriteString(inputValueStyle.Render(m.config.Styling) + "\n")

	// UI Library (if set)
	if m.config.UILibrary != "" && m.config.UILibrary != models.UILibraryNone {
		b.WriteString(compactLabelStyle.Render("UI Library:"))
		b.WriteString(inputValueStyle.Render(m.config.UILibrary) + "\n")
	}

	// Testing (if set)
	if m.config.Testing != "" && m.config.Testing != models.TestingNone {
		b.WriteString(compactLabelStyle.Render("Testing:"))
		b.WriteString(inputValueStyle.Render(m.config.Testing) + "\n")
	}

	// State Management (if set)
	if m.config.StateManagement != "" && m.config.StateManagement != models.StateNone {
		b.WriteString(compactLabelStyle.Render("State Mgmt:"))
		b.WriteString(inputValueStyle.Render(m.config.StateManagement) + "\n")
	}

	// Data Fetching (if set)
	if m.config.DataFetching != "" && m.config.DataFetching != models.DataNone {
		b.WriteString(compactLabelStyle.Render("Data Fetch:"))
		b.WriteString(inputValueStyle.Render(m.config.DataFetching) + "\n")
	}

	// Structure
	if m.config.Structure != "" {
		b.WriteString(compactLabelStyle.Render("Structure:"))
		b.WriteString(inputValueStyle.Render(m.config.Structure) + "\n")
	}

	return b.String()
}

// getEnhancedConfig returns an enhanced configuration display with sections and icons
func (m Model) getEnhancedConfig() string {
	var b strings.Builder

	// Section title style
	sectionStyle := lipgloss.NewStyle().
		Foreground(colorMoltenGold).
		Bold(true).
		Underline(true)

	// Label style with icon support
	labelStyle := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Width(18)

	// Value style with highlighting
	valueStyle := lipgloss.NewStyle().
		Foreground(colorSilverSheen).
		Bold(true)

	// === FOUNDATION SECTION ===
	b.WriteString(sectionStyle.Render("FOUNDATION") + "\n")

	// Project name
	b.WriteString("  " + labelStyle.Render("Project Name:"))
	b.WriteString(valueStyle.Render(m.config.ProjectName) + "\n")

	// Framework with icon
	frameworkIcon := GetFrameworkForgeIcon(m.config.Framework)
	b.WriteString("  " + labelStyle.Render("Framework:"))
	b.WriteString(valueStyle.Render(frameworkIcon+" "+m.config.Framework) + "\n")

	// Language
	b.WriteString("  " + labelStyle.Render("Language:"))
	b.WriteString(valueStyle.Render(m.config.Language) + "\n\n")

	// === TOOLING SECTION ===
	b.WriteString(sectionStyle.Render("TOOLING") + "\n")

	// Package Manager with icon
	pmIcon := GetPackageManagerForgeIcon(m.config.PackageManager)
	b.WriteString("  " + labelStyle.Render("Package Manager:"))
	b.WriteString(valueStyle.Render(pmIcon+" "+m.config.PackageManager) + "\n")

	// Styling with icon
	stylingIcon := GetStylingForgeIcon(m.config.Styling)
	b.WriteString("  " + labelStyle.Render("Styling:"))
	b.WriteString(valueStyle.Render(stylingIcon+m.config.Styling) + "\n")

	// UI Library (if set) with icon
	if m.config.UILibrary != "" && m.config.UILibrary != models.UILibraryNone {
		uiIcon := GetUILibraryForgeIcon(m.config.UILibrary)
		b.WriteString("  " + labelStyle.Render("UI Library:"))
		b.WriteString(valueStyle.Render(uiIcon+m.config.UILibrary) + "\n")
	}

	b.WriteString("\n")

	// === FEATURES SECTION (only if any features are configured) ===
	hasFeatures := (m.config.Routing != "" && m.config.Routing != models.RoutingNone) ||
		(m.config.Testing != "" && m.config.Testing != models.TestingNone) ||
		(m.config.StateManagement != "" && m.config.StateManagement != models.StateNone) ||
		(m.config.FormManagement != "" && m.config.FormManagement != models.FormNone) ||
		(m.config.DataFetching != "" && m.config.DataFetching != models.DataNone)

	if hasFeatures {
		b.WriteString(sectionStyle.Render("FEATURES") + "\n")

		// Routing (if set)
		if m.config.Routing != "" && m.config.Routing != models.RoutingNone {
			b.WriteString("  " + labelStyle.Render("Routing:"))
			b.WriteString(valueStyle.Render(m.config.Routing) + "\n")
		}

		// Testing (if set) with icon
		if m.config.Testing != "" && m.config.Testing != models.TestingNone {
			testIcon := GetTestingForgeIcon(m.config.Testing)
			b.WriteString("  " + labelStyle.Render("Testing:"))
			b.WriteString(valueStyle.Render(testIcon+m.config.Testing) + "\n")
		}

		// State Management (if set) with icon
		if m.config.StateManagement != "" && m.config.StateManagement != models.StateNone {
			stateIcon := GetStateForgeIcon(m.config.StateManagement)
			b.WriteString("  " + labelStyle.Render("State Management:"))
			b.WriteString(valueStyle.Render(stateIcon+m.config.StateManagement) + "\n")
		}

		// Form Management (if set)
		if m.config.FormManagement != "" && m.config.FormManagement != models.FormNone {
			b.WriteString("  " + labelStyle.Render("Form Management:"))
			b.WriteString(valueStyle.Render(m.config.FormManagement) + "\n")
		}

		// Data Fetching (if set)
		if m.config.DataFetching != "" && m.config.DataFetching != models.DataNone {
			b.WriteString("  " + labelStyle.Render("Data Fetching:"))
			b.WriteString(valueStyle.Render(m.config.DataFetching) + "\n")
		}

		b.WriteString("\n")
	}

	// === FINISHING SECTION (only if configured) ===
	hasFinishing := (m.config.Animation != "" && m.config.Animation != models.AnimationNone) ||
		(m.config.Icons != "" && m.config.Icons != models.IconsNone) ||
		(m.config.I18n != "" && m.config.I18n != models.I18nNone) ||
		m.config.Structure != ""

	if hasFinishing {
		b.WriteString(sectionStyle.Render("FINISHING") + "\n")

		// Animation (if set)
		if m.config.Animation != "" && m.config.Animation != models.AnimationNone {
			b.WriteString("  " + labelStyle.Render("Animation:"))
			b.WriteString(valueStyle.Render(m.config.Animation) + "\n")
		}

		// Icons (if set)
		if m.config.Icons != "" && m.config.Icons != models.IconsNone {
			b.WriteString("  " + labelStyle.Render("Icons:"))
			b.WriteString(valueStyle.Render(m.config.Icons) + "\n")
		}

		// i18n (if set)
		if m.config.I18n != "" && m.config.I18n != models.I18nNone {
			b.WriteString("  " + labelStyle.Render("Internationalization:"))
			b.WriteString(valueStyle.Render(m.config.I18n) + "\n")
		}

		// Structure
		if m.config.Structure != "" {
			b.WriteString("  " + labelStyle.Render("Project Structure:"))
			b.WriteString(valueStyle.Render(m.config.Structure) + "\n")
		}
	}

	return b.String()
}

// viewConfirmForge renders a confirmation dialog before starting the forge
func (m Model) viewConfirmForge() string {
	var b strings.Builder
	contentWidth := m.adaptiveBox - 4

	// Confirmation header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorFurnaceOrange).
		Background(colorIronDark).
		Padding(0, 2).
		Width(contentWidth).
		Align(lipgloss.Center).
		Render("CONFIRM ACTION")
	b.WriteString(header + "\n\n")

	// Question
	questionStyle := lipgloss.NewStyle().
		Foreground(colorSilverSheen).
		Bold(true).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(questionStyle.Render("Are you ready to start forging?") + "\n\n")

	// Information
	infoStyle := lipgloss.NewStyle().
		Foreground(colorDraftPencil).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(infoStyle.Render("This will generate your project with the selected configuration.") + "\n\n")

	// Divider
	divider := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render(strings.Repeat("─", contentWidth))
	b.WriteString(divider + "\n\n")

	// Options
	yesOption := lipgloss.NewStyle().
		Foreground(colorTemperedGreen).
		Bold(true).
		Render("[Y] Yes, start forging")

	noOption := lipgloss.NewStyle().
		Foreground(colorCrackedRed).
		Render("[N] No, go back")

	options := yesOption + "    " + noOption
	optionsCentered := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(options)
	b.WriteString(optionsCentered + "\n")

	return lipgloss.NewStyle().
		Padding(2, 0).
		Render(b.String())
}

// viewConfirmBack renders a confirmation dialog before going back to edit
func (m Model) viewConfirmBack() string {
	var b strings.Builder
	contentWidth := m.adaptiveBox - 4

	// Confirmation header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorFlameYellow).
		Background(colorIronDark).
		Padding(0, 2).
		Width(contentWidth).
		Align(lipgloss.Center).
		Render("CONFIRM ACTION")
	b.WriteString(header + "\n\n")

	// Question
	questionStyle := lipgloss.NewStyle().
		Foreground(colorSilverSheen).
		Bold(true).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(questionStyle.Render("Go back to edit configuration?") + "\n\n")

	// Warning
	warningStyle := lipgloss.NewStyle().
		Foreground(colorDraftPencil).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(warningStyle.Render("You'll need to review your changes again.") + "\n\n")

	// Divider
	divider := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render(strings.Repeat("─", contentWidth))
	b.WriteString(divider + "\n\n")

	// Options
	yesOption := lipgloss.NewStyle().
		Foreground(colorSteelBlue).
		Bold(true).
		Render("[Y] Yes, go back")

	noOption := lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render("[N] No, stay here")

	options := yesOption + "    " + noOption
	optionsCentered := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(options)
	b.WriteString(optionsCentered + "\n")

	return lipgloss.NewStyle().
		Padding(2, 0).
		Render(b.String())
}
