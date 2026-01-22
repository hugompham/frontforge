package tui

import "frontforge/internal/models"

// ============================================================================
// FORM PHASES - Logical grouping of form questions
// ============================================================================

// FormPhase represents a logical section of the configuration form
type FormPhase int

const (
	PhaseFoundation FormPhase = iota // Project basics: name, framework, language
	PhaseTooling                     // Build tools: package manager, styling
	PhaseFeatures                    // App features: routing, state, data fetching
	PhaseFinishing                   // Final touches: utils, i18n, structure
)

// String returns the display name for a phase
func (p FormPhase) String() string {
	switch p {
	case PhaseFoundation:
		return "FOUNDATION"
	case PhaseTooling:
		return "TOOLING"
	case PhaseFeatures:
		return "FEATURES"
	case PhaseFinishing:
		return "FINISHING"
	default:
		return "UNKNOWN"
	}
}

// Icon returns the ASCII icon for a phase
// NOTE: Phase icons removed from headers for cleaner design
// Only status icons ([+], [X]) should be used in headers
func (p FormPhase) Icon() string {
	return "" // Icons removed from phase headers
}

// Description returns a brief description of the phase
func (p FormPhase) Description() string {
	switch p {
	case PhaseFoundation:
		return "Core project structure and framework"
	case PhaseTooling:
		return "Development and styling tools"
	case PhaseFeatures:
		return "Application functionality and libraries"
	case PhaseFinishing:
		return "Utilities and project organization"
	default:
		return "Unknown phase"
	}
}

// ============================================================================
// PROGRESS TRACKER - Track progress through form questions
// ============================================================================

// ProgressTracker keeps track of which questions have been answered
type ProgressTracker struct {
	TotalSteps   int       // Total number of questions
	CurrentStep  int       // Current question index (1-based)
	CurrentPhase FormPhase // Current logical phase
	Completed    []bool    // Completion status of each step
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(totalSteps int) *ProgressTracker {
	return &ProgressTracker{
		TotalSteps:   totalSteps,
		CurrentStep:  1,
		CurrentPhase: PhaseFoundation,
		Completed:    make([]bool, totalSteps),
	}
}

// AdvanceStep moves to the next step
func (pt *ProgressTracker) AdvanceStep() {
	if pt.CurrentStep < pt.TotalSteps {
		pt.Completed[pt.CurrentStep-1] = true
		pt.CurrentStep++
		pt.updatePhase()
	}
}

// GoBack returns to the previous step
func (pt *ProgressTracker) GoBack() {
	if pt.CurrentStep > 1 {
		pt.CurrentStep--
		pt.updatePhase()
	}
}

// JumpToStep jumps to a specific step (for editing)
func (pt *ProgressTracker) JumpToStep(step int) {
	if step > 0 && step <= pt.TotalSteps {
		pt.CurrentStep = step
		pt.updatePhase()
	}
}

// PercentComplete returns the completion percentage (0-100)
func (pt *ProgressTracker) PercentComplete() int {
	if pt.TotalSteps == 0 {
		return 0
	}
	return (pt.CurrentStep - 1) * 100 / pt.TotalSteps
}

// IsComplete returns true if all steps are completed
func (pt *ProgressTracker) IsComplete() bool {
	return pt.CurrentStep > pt.TotalSteps
}

// updatePhase updates the current phase based on the current step
// This mapping assumes:
// - Steps 1-3: Foundation (setup mode, project name, language)
// - Steps 4-7: Tooling (framework, package manager, styling, UI library)
// - Steps 8-12: Features (routing, testing, state, form, data fetching)
// - Steps 13+: Finishing (animation, icons, viz, utils, i18n, structure)
func (pt *ProgressTracker) updatePhase() {
	switch {
	case pt.CurrentStep <= 3:
		pt.CurrentPhase = PhaseFoundation
	case pt.CurrentStep <= 7:
		pt.CurrentPhase = PhaseTooling
	case pt.CurrentStep <= 12:
		pt.CurrentPhase = PhaseFeatures
	default:
		pt.CurrentPhase = PhaseFinishing
	}
}

// ============================================================================
// QUESTION METADATA - Information about each form question
// ============================================================================

// QuestionInfo holds metadata about a single form question
type QuestionInfo struct {
	ID          int       // Question number (1-based)
	Phase       FormPhase // Which phase this belongs to
	Title       string    // Question title
	Description string    // Optional help text
	Required    bool      // Whether this question is required
	Conditional bool      // Whether this question depends on previous answers
}

// QuestionCatalog defines all questions in the form
// This helps with navigation and progress tracking
var QuestionCatalog = []QuestionInfo{
	// Foundation Phase
	{ID: 1, Phase: PhaseFoundation, Title: "Setup mode", Description: "Quick preset or custom configuration", Required: true},
	{ID: 2, Phase: PhaseFoundation, Title: "Project name", Description: "Name of your project directory", Required: true},
	{ID: 3, Phase: PhaseFoundation, Title: "Language", Description: "TypeScript or JavaScript", Required: true, Conditional: true},

	// Tooling Phase
	{ID: 4, Phase: PhaseTooling, Title: "Framework", Description: "Your frontend framework", Required: true, Conditional: true},
	{ID: 5, Phase: PhaseTooling, Title: "Package manager", Description: "npm, yarn, pnpm, or bun", Required: true, Conditional: true},
	{ID: 6, Phase: PhaseTooling, Title: "Styling", Description: "CSS framework or methodology", Required: true, Conditional: true},
	{ID: 7, Phase: PhaseTooling, Title: "UI library", Description: "Component library (framework-specific)", Required: false, Conditional: true},

	// Features Phase
	{ID: 8, Phase: PhaseFeatures, Title: "Routing", Description: "Client-side routing solution", Required: false, Conditional: true},
	{ID: 9, Phase: PhaseFeatures, Title: "Testing", Description: "Testing framework", Required: false, Conditional: true},
	{ID: 10, Phase: PhaseFeatures, Title: "State management", Description: "Global state solution", Required: false, Conditional: true},
	{ID: 11, Phase: PhaseFeatures, Title: "Form management", Description: "Form handling library", Required: false, Conditional: true},
	{ID: 12, Phase: PhaseFeatures, Title: "Data fetching", Description: "API communication approach", Required: false, Conditional: true},

	// Finishing Phase
	{ID: 13, Phase: PhaseFinishing, Title: "Animation", Description: "Animation library", Required: false, Conditional: true},
	{ID: 14, Phase: PhaseFinishing, Title: "Icons", Description: "Icon set", Required: false, Conditional: true},
	{ID: 15, Phase: PhaseFinishing, Title: "Data visualization", Description: "Charting library", Required: false, Conditional: true},
	{ID: 16, Phase: PhaseFinishing, Title: "Utilities", Description: "Utility libraries (dates, etc.)", Required: false, Conditional: true},
	{ID: 17, Phase: PhaseFinishing, Title: "Internationalization", Description: "i18n support", Required: false, Conditional: true},
	{ID: 18, Phase: PhaseFinishing, Title: "Project structure", Description: "Folder organization pattern", Required: true, Conditional: true},
}

// GetQuestionInfo returns metadata for a specific question ID
func GetQuestionInfo(id int) *QuestionInfo {
	for i := range QuestionCatalog {
		if QuestionCatalog[i].ID == id {
			return &QuestionCatalog[i]
		}
	}
	return nil
}

// GetQuestionsForPhase returns all questions in a specific phase
func GetQuestionsForPhase(phase FormPhase) []QuestionInfo {
	var questions []QuestionInfo
	for _, q := range QuestionCatalog {
		if q.Phase == phase {
			questions = append(questions, q)
		}
	}
	return questions
}

// ============================================================================
// CONFIGURATION SUMMARY - For review screen
// ============================================================================

// ConfigSection represents a group of related configuration options
type ConfigSection struct {
	Title   string
	Icon    string
	Entries []ConfigEntry
}

// ConfigEntry is a single configuration item
type ConfigEntry struct {
	Label string
	Value string
	Icon  string
}

// BuildConfigSummary creates a structured summary of the configuration
// Useful for the review screen
func BuildConfigSummary(config models.Config) []ConfigSection {
	sections := []ConfigSection{
		{
			Title: "FOUNDATION",
			Icon:  "", // Removed for clean design
			Entries: []ConfigEntry{
				{Label: "Project Name", Value: config.ProjectName},
				{Label: "Language", Value: config.Language},
				{Label: "Framework", Value: config.Framework, Icon: GetFrameworkForgeIcon(config.Framework)},
			},
		},
		{
			Title: "TOOLING",
			Icon:  "", // Removed for clean design
			Entries: []ConfigEntry{
				{Label: "Package Manager", Value: config.PackageManager, Icon: GetPackageManagerForgeIcon(config.PackageManager)},
				{Label: "Styling", Value: config.Styling, Icon: GetStylingForgeIcon(config.Styling)},
			},
		},
	}

	// Add UI library if set
	if config.UILibrary != "" && config.UILibrary != models.UILibraryNone {
		sections[1].Entries = append(sections[1].Entries, ConfigEntry{
			Label: "UI Library",
			Value: config.UILibrary,
			Icon:  GetUILibraryForgeIcon(config.UILibrary),
		})
	}

	// Features section
	featureSection := ConfigSection{
		Title:   "FEATURES",
		Icon:    "", // Removed for clean design
		Entries: []ConfigEntry{},
	}

	if config.Routing != "" && config.Routing != models.RoutingNone {
		featureSection.Entries = append(featureSection.Entries, ConfigEntry{
			Label: "Routing",
			Value: config.Routing,
		})
	}

	if config.StateManagement != "" && config.StateManagement != models.StateNone {
		featureSection.Entries = append(featureSection.Entries, ConfigEntry{
			Label: "State Management",
			Value: config.StateManagement,
			Icon:  GetStateForgeIcon(config.StateManagement),
		})
	}

	if config.DataFetching != "" && config.DataFetching != models.DataNone {
		featureSection.Entries = append(featureSection.Entries, ConfigEntry{
			Label: "Data Fetching",
			Value: config.DataFetching,
		})
	}

	if config.Testing != "" && config.Testing != models.TestingNone {
		featureSection.Entries = append(featureSection.Entries, ConfigEntry{
			Label: "Testing",
			Value: config.Testing,
			Icon:  GetTestingForgeIcon(config.Testing),
		})
	}

	if len(featureSection.Entries) > 0 {
		sections = append(sections, featureSection)
	}

	// Finishing section
	finishingSection := ConfigSection{
		Title:   "FINISHING",
		Icon:    "", // Removed for clean design
		Entries: []ConfigEntry{},
	}

	if config.Structure != "" {
		finishingSection.Entries = append(finishingSection.Entries, ConfigEntry{
			Label: "Project Structure",
			Value: config.Structure,
		})
	}

	if config.I18n != "" && config.I18n != models.I18nNone {
		finishingSection.Entries = append(finishingSection.Entries, ConfigEntry{
			Label: "Internationalization",
			Value: config.I18n,
		})
	}

	if len(finishingSection.Entries) > 0 {
		sections = append(sections, finishingSection)
	}

	return sections
}
