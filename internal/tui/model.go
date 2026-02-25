// Package tui implements the interactive terminal user interface.
//
// This package uses the Bubble Tea framework to provide an interactive,
// keyboard-driven interface for project configuration.
//
// State machine flow:
//
//	Welcome → Form → Confirm → PreflightChecks → Forging → Complete
//
// The TUI follows the Elm Architecture pattern:
//   - Model: Holds application state and configuration
//   - Update: Handles messages and state transitions
//   - View: Renders the current state to terminal
//
// Key components:
//   - Huh forms for structured user input
//   - Spinner for progress indication
//   - Progress tracker for generation steps
//   - Minimal animations (< 300ms) for feedback only
//
// The Model struct is organized into sub-structs for clarity:
//   - formState: Form field values
//   - layout: Terminal dimensions
//   - anim: Simple animation counters
package tui

import (
	"fmt"
	"frontforge/internal/models"
	"frontforge/internal/preflight"
	"frontforge/internal/tui/state"
	"path/filepath"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// State represents the current screen/state in the TUI
type State int

const (
	StateWelcome         State = iota // Welcome screen with forge introduction
	StateTerminalWarn                 // Terminal too narrow warning
	StateBlueprint                    // Form/planning phase (was StateForm)
	StateReview                       // Review configuration before forging
	StateConfirmForge                 // Confirmation dialog before starting forge
	StateConfirmBack                  // Confirmation dialog before going back
	StatePreflightChecks              // Pre-flight validation checks
	StateForging                      // Generation in progress (was StateGenerating)
	StateFinished                     // Success screen (was StateSuccess)
	StateCracked                      // Error screen (was StateError)

	// Legacy state aliases for compatibility
	StateForm       = StateBlueprint
	StateGenerating = StateForging
	StateSuccess    = StateFinished
	StateError      = StateCracked
)

// Model holds the Bubbletea application state
type Model struct {
	// Core state
	currentState  State
	previousState State // Store state before showing terminal warning
	config        models.Config
	err           error

	// Sub-states for better organization
	formState state.FormState
	layout    state.LayoutState
	anim      state.AnimationState

	// Components
	form     *huh.Form
	spinner  spinner.Model
	progress *ProgressTracker

	// Generation state
	generationComplete bool
	preflightResults   *preflight.PreflightResults
}

// Init initializes the Bubbletea model
func (m Model) Init() tea.Cmd {
	// Only initialize spinner, not form yet (form initializes when entering StateBlueprint)
	return m.spinner.Tick
}

// NewModel creates a new Bubbletea model with Huh forms
func NewModel() Model {
	// Create custom forging spinner
	s := spinner.New()
	s.Spinner = ForgingSpinner()
	s.Style = forgeSpinnerStyle

	m := Model{
		currentState:  StateWelcome,
		previousState: StateWelcome,
		config:        models.Config{},
		formState:     state.NewFormState(),
		layout:        state.NewLayoutState(),
		anim:          state.NewAnimationState(),
		spinner:       s,
		progress:      NewProgressTracker(len(QuestionCatalog)),
	}

	m.form = m.createForm()
	return m
}

// NewModelWithPath creates a new Bubbletea model with specified project path
func NewModelWithPath(absPath string, userPath string) Model {
	m := NewModel()

	// Handle three cases:
	// 1. No path provided (absPath == "") - defer path resolution until project name is entered
	// 2. Current directory (userPath == ".") - use current directory
	// 3. Specific folder (userPath != "" && != ".") - use specified folder name

	if absPath == "" {
		// No path specified - will be set later based on project name
		m.config.ProjectPath = ""
		// Keep default project name "my-app"
	} else {
		// Path was explicitly provided via -path flag
		m.config.ProjectPath = absPath

		if userPath == "." {
			// Using current directory - use directory name as default project name
			dirName := filepath.Base(absPath)
			m.formState.ProjectName = dirName
		} else {
			// Using specific folder - use folder name as default project name
			m.formState.ProjectName = userPath
		}
	}

	return m
}

// createForm builds the Huh form with all questions
func (m *Model) createForm() *huh.Form {
	form := huh.NewForm(
		// Group 1: Setup mode
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Setup mode").
				Options(
					huh.NewOption("Quick (use recommended defaults)", string(models.SetupModeQuick)),
					huh.NewOption("Custom (answer each question)", string(models.SetupModeCustom)),
				).
				Value(&m.formState.SetupMode),
		),

		// Group 2: Project name
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Placeholder("my-app").
				Value(&m.formState.ProjectName).
				Validate(func(s string) error {
					if !validateProjectName(s) {
						return fmt.Errorf("project name must contain only letters, numbers, hyphens, and underscores")
					}
					return nil
				}),
		),

		// Group 3: Language (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select your language").
				Options(
					huh.NewOption("TypeScript (recommended)", models.LangTypeScript),
					huh.NewOption("JavaScript", models.LangJavaScript),
				).
				Value(&m.formState.Language),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick)
		}),

		// Group 4: Framework (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a framework").
				Options(
					huh.NewOption("React", models.FrameworkReact),
					huh.NewOption("Vue", models.FrameworkVue),
					huh.NewOption("Angular (standalone)", models.FrameworkAngular),
					huh.NewOption("Svelte", models.FrameworkSvelte),
					huh.NewOption("Solid", models.FrameworkSolid),
					huh.NewOption("Vanilla (no framework)", models.FrameworkVanilla),
					// Meta-frameworks (shell out to upstream CLIs)
					huh.NewOption("Next.js (React)", models.FrameworkNextJS),
					huh.NewOption("Astro (content-focused)", models.FrameworkAstro),
					huh.NewOption("SvelteKit (Svelte)", models.FrameworkSvelteKit),
				).
				Value(&m.formState.Framework),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick)
		}),

		// Group 5: Package manager (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Package manager").
				Options(
					huh.NewOption("npm", models.PackageManagerNpm),
					huh.NewOption("yarn", models.PackageManagerYarn),
					huh.NewOption("pnpm", models.PackageManagerPnpm),
					huh.NewOption("bun", models.PackageManagerBun),
				).
				Value(&m.formState.PackageManager),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick)
		}),

		// Group 6: Styling for Vite-based frameworks (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How would you like to style your app?").
				Options(
					huh.NewOption("Tailwind CSS", models.StylingTailwind),
					huh.NewOption("Bootstrap", models.StylingBootstrap),
					huh.NewOption("CSS Modules", models.StylingCSSModules),
					huh.NewOption("Sass/SCSS", models.StylingSass),
					huh.NewOption("Styled Components", models.StylingStyled),
					huh.NewOption("Vanilla CSS", models.StylingVanilla),
				).
				Value(&m.formState.Styling),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.isMetaSelected()
		}),

		// Group 6-meta: Styling for meta-frameworks (filtered by OptionMatrix)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("How would you like to style your app?").
				Options(
					huh.NewOption("Tailwind CSS", models.StylingTailwind),
					huh.NewOption("CSS Modules", models.StylingCSSModules),
					huh.NewOption("Sass/SCSS", models.StylingSass),
					huh.NewOption("Vanilla CSS", models.StylingVanilla),
				).
				Value(&m.formState.Styling),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || !m.isMetaSelected()
		}),

		// Group 6a: UI Component Library for React (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("UI Component Library").
				Options(
					huh.NewOption("Shadcn/ui (recommended)", models.UILibraryShadcn),
					huh.NewOption("Material-UI (MUI)", models.UILibraryMUI),
					huh.NewOption("Chakra UI", models.UILibraryChakra),
					huh.NewOption("Ant Design", models.UILibraryAntD),
					huh.NewOption("Headless UI", models.UILibraryHeadless),
					huh.NewOption("None", models.UILibraryNone),
				).
				Value(&m.formState.UILibrary),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkReact || m.isMetaSelected()
		}),

		// Group 6b: UI Component Library for Vue (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("UI Component Library").
				Options(
					huh.NewOption("Vuetify", models.UILibraryVuetify),
					huh.NewOption("PrimeVue", models.UILibraryPrimeVue),
					huh.NewOption("Element Plus", models.UILibraryElementUI),
					huh.NewOption("Naive UI", models.UILibraryNaiveUI),
					huh.NewOption("None", models.UILibraryNone),
				).
				Value(&m.formState.UILibrary),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkVue
		}),

		// Group 6c: UI Component Library for Angular (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("UI Component Library").
				Options(
					huh.NewOption("Angular Material", models.UILibraryAngularMaterial),
					huh.NewOption("PrimeNG", models.UILibraryPrimeNG),
					huh.NewOption("NG-ZORRO", models.UILibraryNGZorro),
					huh.NewOption("None", models.UILibraryNone),
				).
				Value(&m.formState.UILibrary),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkAngular
		}),

		// Group 7: Routing for React (only shown in custom mode, not meta-frameworks)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Routing solution").
				Options(
					huh.NewOption("React Router", models.RoutingReactRouter),
					huh.NewOption("TanStack Router", models.RoutingTanStackRouter),
					huh.NewOption("File-based routing", models.RoutingFileBased),
					huh.NewOption("None (single page)", models.RoutingNone),
				).
				Value(&m.formState.Routing),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkReact || m.isMetaSelected()
		}),

		// Group 7b: Routing for Svelte (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Routing solution").
				Options(
					huh.NewOption("SvelteKit (built-in)", models.RoutingSvelteKit),
					huh.NewOption("None (single page)", models.RoutingNone),
				).
				Value(&m.formState.Routing),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkSvelte
		}),

		// Group 7c: Routing for Solid (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Routing solution").
				Options(
					huh.NewOption("Solid Router", models.RoutingSolidRouter),
					huh.NewOption("None (single page)", models.RoutingNone),
				).
				Value(&m.formState.Routing),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkSolid
		}),

		// Group 8: Testing for Vite-based (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Testing framework").
				Options(
					huh.NewOption("Vitest (fast, Vite-native)", models.TestingVitest),
					huh.NewOption("Jest", models.TestingJest),
					huh.NewOption("None (set up later)", models.TestingNone),
				).
				Value(&m.formState.Testing),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.isMetaSelected()
		}),

		// Group 8-meta: Testing for Next.js (Vitest + Jest + None)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Testing framework").
				Options(
					huh.NewOption("Vitest (fast, Vite-native)", models.TestingVitest),
					huh.NewOption("Jest", models.TestingJest),
					huh.NewOption("None (set up later)", models.TestingNone),
				).
				Value(&m.formState.Testing),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkNextJS
		}),

		// Group 8-meta: Testing for SvelteKit (Vitest + Playwright + None)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Testing framework").
				Options(
					huh.NewOption("Vitest (fast, Vite-native)", models.TestingVitest),
					huh.NewOption("Playwright (E2E)", models.TestingPlaywright),
					huh.NewOption("None (set up later)", models.TestingNone),
				).
				Value(&m.formState.Testing),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkSvelteKit
		}),

		// Group 8-meta: Testing for Astro (Vitest + None)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Testing framework").
				Options(
					huh.NewOption("Vitest (fast, Vite-native)", models.TestingVitest),
					huh.NewOption("None (set up later)", models.TestingNone),
				).
				Value(&m.formState.Testing),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkAstro
		}),

		// Group 9: State management for React (only shown in custom mode, not meta)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("Zustand (lightweight)", models.StateZustand),
					huh.NewOption("Redux Toolkit", models.StateReduxToolkit),
					huh.NewOption("Context API only", models.StateContextAPI),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.formState.StateManagement),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkReact || m.isMetaSelected()
		}),

		// Group 9-meta: State management for Next.js
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("Zustand (lightweight)", models.StateZustand),
					huh.NewOption("Redux Toolkit", models.StateReduxToolkit),
					huh.NewOption("Context API only", models.StateContextAPI),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.formState.StateManagement),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkNextJS
		}),

		// Group 9-meta: State management for SvelteKit
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("Svelte Stores (built-in)", models.StateSvelteStores),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.formState.StateManagement),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkSvelteKit
		}),

		// Group 9b: State management for Vue (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("Pinia (recommended)", models.StatePinia),
					huh.NewOption("Vuex", models.StateVuex),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.formState.StateManagement),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkVue
		}),

		// Group 9c: State management for Svelte (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("Svelte Stores (built-in)", models.StateSvelteStores),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.formState.StateManagement),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkSvelte
		}),

		// Group 9d: State management for Solid (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("Solid Stores (built-in)", models.StateSolidStores),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.formState.StateManagement),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkSolid
		}),

		// Group 9e: State management for Angular (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("NgRx", models.StateNgRx),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.formState.StateManagement),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkAngular
		}),

		// Group 9f: Form Management for React (only shown in custom mode, not meta)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Form Management").
				Options(
					huh.NewOption("React Hook Form (recommended)", models.FormReactHookForm),
					huh.NewOption("Formik", models.FormFormik),
					huh.NewOption("TanStack Form", models.FormTanStackForm),
					huh.NewOption("None (native)", models.FormNone),
				).
				Value(&m.formState.FormManagement),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkReact || m.isMetaSelected()
		}),

		// Group 9g: Form Management for Vue (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Form Management").
				Options(
					huh.NewOption("VeeValidate", models.FormVeeValidate),
					huh.NewOption("None (native)", models.FormNone),
				).
				Value(&m.formState.FormManagement),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkVue
		}),

		// Group 10: Data fetching for Vite-based (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Data fetching approach").
				Options(
					huh.NewOption("TanStack Query", models.DataTanStackQuery),
					huh.NewOption("Fetch API", models.DataFetchAPI),
					huh.NewOption("Axios", models.DataAxios),
					huh.NewOption("SWR", models.DataSWR),
					huh.NewOption("None (no API calls)", models.DataNone),
				).
				Value(&m.formState.DataFetching),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.isMetaSelected()
		}),

		// Group 10-meta: Data fetching for Next.js (TanStack + SWR + Axios + Fetch + None)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Data fetching approach").
				Options(
					huh.NewOption("TanStack Query", models.DataTanStackQuery),
					huh.NewOption("SWR", models.DataSWR),
					huh.NewOption("Axios", models.DataAxios),
					huh.NewOption("Fetch API", models.DataFetchAPI),
					huh.NewOption("None (no API calls)", models.DataNone),
				).
				Value(&m.formState.DataFetching),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkNextJS
		}),

		// Group 10-meta: Data fetching for SvelteKit (TanStack + Fetch + None)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Data fetching approach").
				Options(
					huh.NewOption("TanStack Query", models.DataTanStackQuery),
					huh.NewOption("Fetch API", models.DataFetchAPI),
					huh.NewOption("None (no API calls)", models.DataNone),
				).
				Value(&m.formState.DataFetching),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.formState.Framework != models.FrameworkSvelteKit
		}),

		// Group 10a: Animation (only shown in custom mode, not meta)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Animation library").
				Options(
					huh.NewOption("Framer Motion (React)", models.AnimationFramerMotion),
					huh.NewOption("GSAP (any framework)", models.AnimationGSAP),
					huh.NewOption("Auto Animate", models.AnimationAutoAnimate),
					huh.NewOption("React Spring", models.AnimationReactSpring),
					huh.NewOption("None", models.AnimationNone),
				).
				Value(&m.formState.Animation),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.isMetaSelected()
		}),

		// Group 10b: Icons (only shown in custom mode, not meta)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Icon library").
				Options(
					huh.NewOption("Heroicons", models.IconsHeroicons),
					huh.NewOption("Lucide", models.IconsLucide),
					huh.NewOption("React Icons", models.IconsReactIcons),
					huh.NewOption("Font Awesome", models.IconsFontAwesome),
					huh.NewOption("None", models.IconsNone),
				).
				Value(&m.formState.Icons),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.isMetaSelected()
		}),

		// Group 10c: Data Visualization (only shown in custom mode, not meta)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Data visualization").
				Options(
					huh.NewOption("Recharts (React)", models.DataVizRecharts),
					huh.NewOption("Chart.js (any framework)", models.DataVizChartJS),
					huh.NewOption("Apache ECharts", models.DataVizECharts),
					huh.NewOption("Nivo (React)", models.DataVizNivo),
					huh.NewOption("None", models.DataVizNone),
				).
				Value(&m.formState.DataViz),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.isMetaSelected()
		}),

		// Group 10d: Utilities (only shown in custom mode, not meta)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Utility library").
				Options(
					huh.NewOption("date-fns (date manipulation)", models.UtilsDateFns),
					huh.NewOption("Day.js (lightweight dates)", models.UtilsDayJS),
					huh.NewOption("Lodash-es (tree-shakeable)", models.UtilsLodash),
					huh.NewOption("None", models.UtilsNone),
				).
				Value(&m.formState.Utilities),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.isMetaSelected()
		}),

		// Group 10e: Internationalization (only shown in custom mode, not meta)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Internationalization (i18n)").
				Options(
					huh.NewOption("react-i18next (React)", models.I18nReactI18next),
					huh.NewOption("vue-i18n (Vue)", models.I18nVueI18n),
					huh.NewOption("None", models.I18nNone),
				).
				Value(&m.formState.I18n),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick) || m.isMetaSelected()
		}),

		// Group 11: Folder structure (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Folder structure").
				Options(
					huh.NewOption("Feature-based (recommended)", models.StructureFeatureBased),
					huh.NewOption("Layer-based", models.StructureLayerBased),
				).
				Value(&m.formState.Structure),
		).WithHideFunc(func() bool {
			return m.formState.SetupMode == string(models.SetupModeQuick)
		}),
		// No confirmation question - review screen handles this
	)

	// Apply custom forge theme to match the rest of the TUI
	form.WithTheme(ForgeTheme())

	return form
}

// isMetaSelected returns true if the current framework selection is a meta-framework
func (m *Model) isMetaSelected() bool {
	return models.IsMetaFramework(m.formState.Framework)
}

// validateProjectName checks if project name is valid
func validateProjectName(name string) bool {
	if name == "" {
		return false
	}
	// Must contain only alphanumeric, hyphens, and underscores
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

// ValidateProjectName is the public version of validateProjectName for testing
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name is required")
	}
	if !validateProjectName(name) {
		return fmt.Errorf("project name must contain only alphanumeric characters, hyphens, and underscores")
	}
	return nil
}

// applyFormDataToConfig converts form field values to config
func (m *Model) applyFormDataToConfig() {
	m.config.ProjectName = m.formState.ProjectName

	// If ProjectPath wasn't set (no -path flag), set it based on project name
	if m.config.ProjectPath == "" {
		cwd, err := filepath.Abs(".")
		if err == nil {
			m.config.ProjectPath = filepath.Join(cwd, m.formState.ProjectName)
		}
	}

	// Apply quick preset if selected
	if m.formState.SetupMode == string(models.SetupModeQuick) {
		preset := models.QuickPreset()
		m.config.Language = preset.Language
		m.config.Framework = preset.Framework
		m.config.PackageManager = preset.PackageManager
		m.config.Styling = preset.Styling
		m.config.Routing = preset.Routing
		m.config.Testing = preset.Testing
		m.config.StateManagement = preset.StateManagement
		m.config.DataFetching = preset.DataFetching
		m.config.Structure = preset.Structure
		return
	}

	// Apply custom configuration
	m.config.Language = m.formState.Language
	m.config.Framework = m.formState.Framework
	m.config.PackageManager = m.formState.PackageManager
	m.config.Styling = m.formState.Styling
	m.config.UILibrary = m.formState.UILibrary
	m.config.Testing = m.formState.Testing
	m.config.FormManagement = m.formState.FormManagement
	m.config.DataFetching = m.formState.DataFetching
	m.config.Animation = m.formState.Animation
	m.config.Icons = m.formState.Icons
	m.config.DataViz = m.formState.DataViz
	m.config.Utilities = m.formState.Utilities
	m.config.I18n = m.formState.I18n
	m.config.Structure = m.formState.Structure

	// Handle framework-specific defaults
	switch m.formState.Framework {
	case models.FrameworkVanilla:
		m.config.Routing = models.RoutingNone
		m.config.StateManagement = models.StateNone
		m.config.UILibrary = models.UILibraryNone
		m.config.FormManagement = models.FormNone
	case models.FrameworkVue:
		m.config.Routing = models.RoutingVueRouter
		m.config.StateManagement = m.formState.StateManagement
	case models.FrameworkAngular:
		m.config.Routing = models.RoutingAngularRouter
		m.config.StateManagement = m.formState.StateManagement
	case models.FrameworkSvelte:
		m.config.Routing = m.formState.Routing
		m.config.StateManagement = m.formState.StateManagement
		m.config.UILibrary = models.UILibraryNone
		m.config.FormManagement = models.FormNone
	case models.FrameworkSolid:
		m.config.Routing = m.formState.Routing
		m.config.StateManagement = m.formState.StateManagement
		m.config.UILibrary = models.UILibraryNone
		m.config.FormManagement = models.FormNone
	case models.FrameworkReact:
		m.config.Routing = m.formState.Routing
		m.config.StateManagement = m.formState.StateManagement
	case models.FrameworkNextJS:
		m.config.Routing = models.RoutingNextJSAppRouter
		m.config.StateManagement = m.formState.StateManagement
		m.config.UILibrary = models.UILibraryNone
		m.config.FormManagement = models.FormNone
	case models.FrameworkAstro:
		m.config.Routing = models.RoutingAstroPages
		m.config.StateManagement = models.StateNone
		m.config.UILibrary = models.UILibraryNone
		m.config.FormManagement = models.FormNone
	case models.FrameworkSvelteKit:
		m.config.Routing = models.RoutingSvelteKit
		m.config.StateManagement = models.StateSvelteStores
		m.config.UILibrary = models.UILibraryNone
		m.config.FormManagement = models.FormNone
	}
}

// Public methods for testing

// GetCurrentState returns the current state
func (m Model) GetCurrentState() State {
	return m.currentState
}

// SetCurrentState sets the current state (for testing)
func (m *Model) SetCurrentState(s State) {
	m.currentState = s
}

// GetFormState returns the form state
func (m Model) GetFormState() state.FormState {
	return m.formState
}

// SetFormState sets the form state (for testing)
func (m *Model) SetFormState(fs state.FormState) {
	m.formState = fs
}

// GetConfig returns the config
func (m Model) GetConfig() models.Config {
	return m.config
}

// ApplyFormDataToConfig is a public version of applyFormDataToConfig for testing
func (m *Model) ApplyFormDataToConfig() {
	m.applyFormDataToConfig()
}
