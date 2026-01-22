package tui

import (
	"fmt"
	"frontforge/internal/models"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================================
// FORGING THEME COLOR PALETTE
// ============================================================================

var (
	// Furnace Heat (Primary) - Warm, energetic colors for active elements
	colorFurnaceOrange = lipgloss.Color("#E85D3B") // Hot metal, primary actions (improved contrast)
	colorEmberRed      = lipgloss.Color("#D64545") // Glowing embers, active borders
	colorMoltenGold    = lipgloss.Color("#FFA500") // Molten metal, progress/selection
	colorFlameYellow   = lipgloss.Color("#F4C430") // Flame tips, highlights (richer saffron)

	// Steel & Metal (Secondary) - Cool, structural colors
	colorSteelBlue   = lipgloss.Color("#4A5F7F") // Cool steel, secondary borders
	colorAnvilGray   = lipgloss.Color("#9E9E9E") // Anvil surface, labels (lighter for better contrast)
	colorIronDark    = lipgloss.Color("#2D2D2D") // Raw iron, backgrounds (darker for depth)
	colorSilverSheen = lipgloss.Color("#BDBDBD") // Polished metal, values (warmer silver)

	// Blueprint (Planning Phase) - Cool, thoughtful colors
	colorBlueprintBlue = lipgloss.Color("#2E5C8A") // Blueprint background
	colorDraftPencil   = lipgloss.Color("#8AA6C2") // Sketch lines, prompts

	// Accent & Feedback
	colorSparkGold = lipgloss.Color("#FFB700") // Sparks flying, code blocks (richer amber)
	colorCharcoal  = lipgloss.Color("#1E1E1E") // Forge charcoal, code backgrounds (truer charcoal)
	colorAshGray   = lipgloss.Color("#B0B0B0") // Muted/disabled text

	// Status Colors
	colorTemperedGreen = lipgloss.Color("#26A69A") // Success (tempered steel, material teal)
	colorCrackedRed    = lipgloss.Color("#D32F2F") // Error (cracked metal, less harsh red)

	// Legacy colors (kept for compatibility during transition)
	colorCyan    = lipgloss.Color("14")
	colorGreen   = lipgloss.Color("10")
	colorRed     = lipgloss.Color("9")
	colorYellow  = lipgloss.Color("11")
	colorGray    = lipgloss.Color("245")
	colorWhite   = lipgloss.Color("15")
	colorBgBlack = lipgloss.Color("0")

	// Box dimensions
	boxWidth    = 76  // Default width for better breathing room and modern terminals
	minBoxWidth = 60  // Minimum width for narrow terminals
	maxBoxWidth = 100 // Maximum width for wide terminals
)

// ============================================================================
// BORDER DEFINITIONS - Industrial, heavy borders for forging theme
// ============================================================================

var (
	// Heavy Forge Border - Thick, industrial feel
	forgeBorder = lipgloss.Border{
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}

	// Anvil Border - Extra thick double lines for heavy emphasis
	anvilBorder = lipgloss.Border{
		Top:         "═",
		Bottom:      "═",
		Left:        "║",
		Right:       "║",
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
	}

	// Blueprint Border - Light, drafting lines for planning phase
	blueprintBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
	}
)

// ============================================================================
// TEXT STYLES - Typography and formatting
// ============================================================================

var (
	// Text styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	labelStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Width(16)

	valueStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	successTextStyle = lipgloss.NewStyle().
				Foreground(colorGreen).
				Bold(true)

	errorTextStyle = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	warningTextStyle = lipgloss.NewStyle().
				Foreground(colorYellow)

	infoTextStyle = lipgloss.NewStyle().
			Foreground(colorCyan)

	// Question styles
	questionStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	choiceStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	selectedChoiceStyle = lipgloss.NewStyle().
				Foreground(colorCyan).
				Bold(true)

	// Spinner style
	spinnerStyle = lipgloss.NewStyle().
			Foreground(colorCyan)

	// Code style
	codeStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Background(colorBgBlack).
			Padding(0, 1)

	// Muted style
	mutedStyle = lipgloss.NewStyle().
			Foreground(colorGray)
)

// ============================================================================
// FORGING THEME TEXT STYLES - Typography hierarchy for the forge
// ============================================================================

var (
	// Headers - Bold, impactful, commands attention
	forgeHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorFurnaceOrange).
				Underline(true).
				MarginBottom(1)

	// Phase Titles - For major sections (BLUEPRINT PHASE, FORGING, etc.)
	phaseTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorMoltenGold).
			Padding(0, 2).
			Background(colorIronDark)

	// Section Headers - For subsections (FOUNDATION, TOOLING, etc.)
	sectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorSilverSheen).
				Underline(true)

	// Prompts - Question text for user input
	promptStyle = lipgloss.NewStyle().
			Foreground(colorDraftPencil).
			Bold(true)

	// Input Values - User-entered or selected values
	inputValueStyle = lipgloss.NewStyle().
			Foreground(colorSilverSheen)
		// Removed italic for better readability

	// Forge Labels - Left-aligned labels for config display
	forgeLabelStyle = lipgloss.NewStyle().
			Foreground(colorAnvilGray).
			Width(20). // Increased from 18 for better alignment
			Align(lipgloss.Left)

	// Help Text - Descriptive, explanatory text
	helpTextStyle = lipgloss.NewStyle().
			Foreground(colorAshGray)
		// Removed italic for better accessibility and readability

	// Recommended Badge - Highlights recommended options
	recommendedBadgeStyle = lipgloss.NewStyle().
				Foreground(colorCharcoal).
				Background(colorSparkGold).
				Padding(0, 1).
				Bold(true)

	// Success Text - Completion messages
	forgeSuccessStyle = lipgloss.NewStyle().
				Foreground(colorTemperedGreen).
				Bold(true)

	// Error Text - Failure messages
	forgeErrorStyle = lipgloss.NewStyle().
			Foreground(colorCrackedRed).
			Bold(true)

	// Warning Text - Caution messages
	forgeWarningStyle = lipgloss.NewStyle().
				Foreground(colorFlameYellow).
				Bold(true)

	// Info Text - Informational messages
	forgeInfoStyle = lipgloss.NewStyle().
			Foreground(colorBlueprintBlue)

	// Muted Text - De-emphasized content
	forgeMutedStyle = lipgloss.NewStyle().
			Foreground(colorAshGray)

	// Spinner Style - Loading indicator
	forgeSpinnerStyle = lipgloss.NewStyle().
				Foreground(colorMoltenGold).
				Bold(true)
)

// ============================================================================
// HELPER FUNCTIONS - Common rendering utilities
// ============================================================================

func RenderConfigRow(label, value string) string {
	return labelStyle.Render(label) + valueStyle.Render(value)
}

func RenderSuccess(msg string) string {
	return successTextStyle.Render(msg)
}

func RenderWarning(msg string) string {
	return warningTextStyle.Render(msg)
}

func RenderInfo(msg string) string {
	return infoTextStyle.Render(msg)
}

// Forging-themed helper functions
func RenderForgeHeader(title string) string {
	return forgeHeaderStyle.Render(title)
}

func RenderPhaseTitle(phase string) string {
	return phaseTitleStyle.Render(fmt.Sprintf("  %s  ", phase))
}

func RenderSectionHeader(section string) string {
	return sectionHeaderStyle.Render(section)
}

func RenderForgeConfigRow(label, value string) string {
	return forgeLabelStyle.Render(label) + inputValueStyle.Render(value)
}

func RenderRecommendedBadge() string {
	return recommendedBadgeStyle.Render(" RECOMMENDED ")
}

// RenderAdaptiveDivider creates a divider with specified width
// Use this for responsive layouts
func RenderAdaptiveDivider(contentWidth int) string {
	return lipgloss.NewStyle().
		Foreground(colorAnvilGray).
		Render(strings.Repeat("─", contentWidth))
}

// RenderProgressBar creates a visual progress bar
// current: current step number, total: total steps, width: bar width in characters
func RenderProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(float64(width) * percentage)
	if filled > width {
		filled = width
	}
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	percentText := fmt.Sprintf("%d%%", int(percentage*100))

	progressStyle := lipgloss.NewStyle().
		Foreground(colorMoltenGold)

	return progressStyle.Render(bar) + "  " + percentText
}

// RenderStepCounter shows current step out of total
func RenderStepCounter(current, total int) string {
	return forgeMutedStyle.Render(fmt.Sprintf("Step %d/%d", current, total))
}

// Framework and library icons (Unicode/Nerd Font compatible)
var frameworkIcons = map[string]string{
	models.FrameworkReact:   "⚛",  // React atom
	models.FrameworkVue:     "V",  // Vue (simple V for compatibility)
	models.FrameworkAngular: "A",  // Angular
	models.FrameworkSvelte:  "S",  // Svelte
	models.FrameworkSolid:   "◆",  // Solid
	models.FrameworkVanilla: "JS", // Vanilla JS
}

var packageManagerIcons = map[string]string{
	models.PackageManagerNpm:  "N", // npm
	models.PackageManagerYarn: "Y", // yarn
	models.PackageManagerPnpm: "P", // pnpm
	models.PackageManagerBun:  "B", // bun
}

var uiLibraryIcons = map[string]string{
	models.UILibraryShadcn:          "S", // Shadcn
	models.UILibraryMUI:             "M", // MUI
	models.UILibraryChakra:          "C", // Chakra
	models.UILibraryAntD:            "A", // Ant Design
	models.UILibraryHeadless:        "H", // Headless
	models.UILibraryVuetify:         "V", // Vuetify
	models.UILibraryPrimeVue:        "P", // PrimeVue
	models.UILibraryElementUI:       "E", // Element
	models.UILibraryNaiveUI:         "N", // Naive UI
	models.UILibraryAngularMaterial: "M", // Angular Material
	models.UILibraryPrimeNG:         "P", // PrimeNG
	models.UILibraryNGZorro:         "Z", // NG-ZORRO
}

var stylingIcons = map[string]string{
	models.StylingTailwind:   "T",  // Tailwind
	models.StylingBootstrap:  "B",  // Bootstrap
	models.StylingSass:       "S",  // Sass
	models.StylingStyled:     "SC", // Styled Components
	models.StylingCSSModules: "M",  // CSS Modules
	models.StylingVanilla:    "C",  // CSS
}

var testingIcons = map[string]string{
	models.TestingVitest: "V", // Vitest
	models.TestingJest:   "J", // Jest
}

var stateIcons = map[string]string{
	models.StateZustand:      "Z", // Zustand
	models.StateReduxToolkit: "R", // Redux
	models.StatePinia:        "P", // Pinia
}

// GetFrameworkIcon returns the icon for a framework
func GetFrameworkIcon(framework string) string {
	if icon, ok := frameworkIcons[framework]; ok {
		return icon
	}
	return "F" // Default folder icon
}

// GetPackageManagerIcon returns the icon for a package manager
func GetPackageManagerIcon(pm string) string {
	if icon, ok := packageManagerIcons[pm]; ok {
		return icon
	}
	return "P" // Default package icon
}

// GetUILibraryIcon returns the icon for a UI library
func GetUILibraryIcon(lib string) string {
	if lib == models.UILibraryNone || lib == "" {
		return ""
	}
	if icon, ok := uiLibraryIcons[lib]; ok {
		return icon + " "
	}
	return "U " // Default UI icon
}

// GetStylingIcon returns the icon for a styling solution
func GetStylingIcon(styling string) string {
	if icon, ok := stylingIcons[styling]; ok {
		return icon + " "
	}
	return "S " // Default style icon
}

// GetTestingIcon returns the icon for testing framework
func GetTestingIcon(testing string) string {
	if testing == models.TestingNone || testing == "" {
		return ""
	}
	if icon, ok := testingIcons[testing]; ok {
		return icon + " "
	}
	return "T " // Default testing icon
}

// GetStateIcon returns the icon for state management
func GetStateIcon(state string) string {
	if state == models.StateNone || state == "" {
		return ""
	}
	if icon, ok := stateIcons[state]; ok {
		return icon + " "
	}
	return "S " // Default state icon
}

// ============================================================================
// RESPONSIVE LAYOUT HELPERS
// ============================================================================

// CalculateBoxWidth returns an appropriate box width based on terminal width
// Ensures content fits well across different terminal sizes
func CalculateBoxWidth(terminalWidth int) int {
	// Edge case: extremely narrow terminal
	if terminalWidth < 40 {
		// Return minimum viable width, even if it might overflow
		return 36 // Enough for basic content
	}

	// Leave padding on both sides (4 chars each side = 8 total)
	availableWidth := terminalWidth - 8

	// Clamp between min and max
	if availableWidth < minBoxWidth {
		return minBoxWidth
	}
	if availableWidth > maxBoxWidth {
		return maxBoxWidth
	}
	return availableWidth
}

// IsTerminalTooNarrow checks if terminal is below minimum usable width
func IsTerminalTooNarrow(terminalWidth int) bool {
	return terminalWidth < 80 // Minimum 80 characters for good UX
}

// GetAdaptiveStyles returns box styles adapted to current terminal width
func GetAdaptiveStyles(terminalWidth int) (welcomeBox, blueprintBox, forgeBox, temperedBox, crackedBox, reviewBox lipgloss.Style) {
	adaptiveWidth := CalculateBoxWidth(terminalWidth)

	welcomeBox = lipgloss.NewStyle().
		Border(anvilBorder).
		BorderForeground(colorFurnaceOrange).
		Width(adaptiveWidth).
		Align(lipgloss.Center).
		Padding(1, 2)

	blueprintBox = lipgloss.NewStyle().
		Border(blueprintBorder).
		BorderForeground(colorBlueprintBlue).
		Width(adaptiveWidth).
		Padding(1, 2)

	forgeBox = lipgloss.NewStyle().
		Border(forgeBorder).
		BorderForeground(colorFurnaceOrange).
		Width(adaptiveWidth).
		Padding(1, 2)

	// SUCCESS uses double border (anvilBorder) to differentiate from error
	temperedBox = lipgloss.NewStyle().
		Border(anvilBorder).
		BorderForeground(colorTemperedGreen).
		Width(adaptiveWidth).
		Padding(1, 2)

	// ERROR uses heavy border (forgeBorder) to differentiate from success
	crackedBox = lipgloss.NewStyle().
		Border(forgeBorder).
		BorderForeground(colorCrackedRed).
		Width(adaptiveWidth).
		Padding(1, 2)

	reviewBox = lipgloss.NewStyle().
		Border(anvilBorder).
		BorderForeground(colorSteelBlue).
		Width(adaptiveWidth).
		Padding(1, 2)

	return
}

// RenderStatusHeader renders a header with both icon and border pattern for accessibility
// Success uses [+] with double border, Error uses [X] with heavy border, Warning uses [!]
func RenderStatusHeader(statusType string, message string) string {
	var icon, colorCode string
	var style lipgloss.Style

	switch statusType {
	case "success":
		icon = "[+]"
		colorCode = "#26A69A" // Tempered green
		// Double equals for success
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCode)).
			Bold(true)
	case "error":
		icon = "[X]"
		colorCode = "#D32F2F" // Cracked red
		// Heavy bars for error
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCode)).
			Bold(true)
	case "warning":
		icon = "[!]"
		colorCode = "#F4C430" // Flame yellow
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCode)).
			Bold(true)
	default:
		icon = "[i]"
		colorCode = "#4A5F7F" // Steel blue
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCode))
	}

	return style.Render(icon + " " + message + " " + icon)
}

// ============================================================================
// HUH FORM THEME - Custom forge theme for Huh forms
// ============================================================================

// ForgeTheme returns a custom Huh theme matching the forge color palette
func ForgeTheme() *huh.Theme {
	t := huh.ThemeBase()

	// Focused field styles - active, being edited
	t.Focused.Title = t.Focused.Title.Foreground(colorMoltenGold).Bold(true)
	t.Focused.Description = t.Focused.Description.Foreground(colorDraftPencil)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(colorFurnaceOrange)
	t.Focused.Option = t.Focused.Option.Foreground(colorSilverSheen)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(colorMoltenGold).Bold(true)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(colorFurnaceOrange)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(colorAnvilGray)
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(colorAshGray)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(colorCharcoal).Background(colorMoltenGold).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(colorAshGray)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(colorCrackedRed)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(colorCrackedRed)

	// Blurred field styles - not currently active
	t.Blurred.Title = t.Blurred.Title.Foreground(colorAnvilGray)
	t.Blurred.Description = t.Blurred.Description.Foreground(colorAshGray)
	t.Blurred.SelectSelector = t.Blurred.SelectSelector.Foreground(colorAshGray)
	t.Blurred.Option = t.Blurred.Option.Foreground(colorAshGray)
	t.Blurred.SelectedOption = t.Blurred.SelectedOption.Foreground(colorAnvilGray)
	t.Blurred.FocusedButton = t.Blurred.FocusedButton.Foreground(colorAshGray).Background(colorIronDark)
	t.Blurred.BlurredButton = t.Blurred.BlurredButton.Foreground(colorAshGray)

	// Help text
	t.Help.Ellipsis = t.Help.Ellipsis.Foreground(colorAshGray)
	t.Help.ShortKey = t.Help.ShortKey.Foreground(colorDraftPencil)
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(colorAnvilGray)
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(colorAshGray)
	t.Help.FullKey = t.Help.FullKey.Foreground(colorMoltenGold)
	t.Help.FullDesc = t.Help.FullDesc.Foreground(colorSilverSheen)
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(colorAnvilGray)

	return t
}
