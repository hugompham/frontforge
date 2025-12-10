package tui

import (
	"frontforge/internal/models"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

// State represents the current screen/state in the TUI
type State int

const (
	StateWelcome         State = iota // Welcome screen with forge introduction
	StateTerminalWarn                  // Terminal too narrow warning
	StateBlueprint                     // Form/planning phase (was StateForm)
	StateReview                        // Review configuration before forging
	StateConfirmForge                  // Confirmation dialog before starting forge
	StateConfirmBack                   // Confirmation dialog before going back
	StateForging                       // Generation in progress (was StateGenerating)
	StateFinished                      // Success screen (was StateSuccess)
	StateCracked                       // Error screen (was StateError)

	// Legacy state aliases for compatibility
	StateForm       = StateBlueprint
	StateGenerating = StateForging
	StateSuccess    = StateFinished
	StateError      = StateCracked
)

// Model holds the Bubbletea application state
type Model struct {
	state   State
	config  models.Config
	err     error
	spinner spinner.Model

	// Huh form
	form *huh.Form

	// Temporary form fields
	setupMode       string
	projectName     string
	language        string
	framework       string
	packageManager  string
	styling         string
	uiLibrary       string
	routing         string
	testing         string
	stateManagement string
	formManagement  string
	dataFetching    string
	animation       string
	icons           string
	dataViz         string
	utilities       string
	i18n            string
	structure       string

	// Generation state
	generationComplete bool

	// Progress tracking
	progress    *ProgressTracker
	tickCount   int    // Animation tick counter
	currentTask string // Current generation task description

	// Terminal dimensions
	width       int
	height      int
	adaptiveBox int // Adaptive box width based on terminal size

	// State management
	previousState State // Store state before showing terminal warning

	// Animation physics - Welcome screen
	flameIntensity     float64          // Current flame intensity (0.0 to 1.0)
	flameVelocity      float64          // Flame velocity for spring physics
	flameSpring        harmonica.Spring // Spring for smooth flame animation
	flameTargetToggle  bool             // Toggle between low and high flame targets

	// Animation physics - Forging screen
	forgingFlameIntensity     float64          // Forging flame intensity (0.0 to 1.0)
	forgingFlameVelocity      float64          // Forging flame velocity
	forgingFlameSpring        harmonica.Spring // Spring for forging flame animation
	forgingFlameTargetToggle  bool             // Toggle for forging flames
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

	// Initialize with safe defaults
	initialWidth := 80
	initialHeight := 24

	m := Model{
		state:           StateWelcome, // Start with welcome screen
		config:          models.Config{},
		spinner:         s,
		width:           initialWidth,
		height:          initialHeight,
		adaptiveBox:     CalculateBoxWidth(initialWidth), // Calculate adaptive width from start
		tickCount:       0,
		currentTask:     "",
		progress:        NewProgressTracker(len(QuestionCatalog)),
		previousState:   StateWelcome,
		setupMode:       string(models.SetupModeCustom),
		projectName:     "my-app",
		language:        models.LangTypeScript,
		framework:       models.FrameworkReact,
		packageManager:  models.PackageManagerNpm,
		styling:         models.StylingTailwind,
		uiLibrary:       models.UILibraryShadcn,
		routing:         models.RoutingReactRouter,
		testing:         models.TestingVitest,
		stateManagement: models.StateZustand,
		formManagement:  models.FormReactHookForm,
		dataFetching:    models.DataTanStackQuery,
		animation:       models.AnimationFramerMotion,
		icons:           models.IconsHeroicons,
		dataViz:         models.DataVizNone,
		utilities:       models.UtilsDateFns,
		i18n:            models.I18nNone,
		structure:       models.StructureFeatureBased,
		// Initialize flame animation with spring physics (Welcome screen)
		// FPS(60): 60 frames per second
		// 0.2: Angular frequency - incredibly slow, barely perceptible breathing
		// 2.0: Damping ratio - massively over-damped (glacial pace, ultra smooth)
		flameSpring:       harmonica.NewSpring(harmonica.FPS(60), 0.2, 2.0),
		flameIntensity:    0.35, // Start at low flame
		flameVelocity:     0.0,
		flameTargetToggle: false, // Start with low target

		// Initialize forging hammer animation (Forging screen)
		// Match the welcome furnace animation speed
		// FPS(60): 60 frames per second
		// 0.2: Angular frequency - same as furnace (incredibly slow, barely perceptible)
		// 2.0: Damping ratio - same as furnace (massively over-damped, glacial pace)
		forgingFlameSpring:       harmonica.NewSpring(harmonica.FPS(60), 0.2, 2.0),
		forgingFlameIntensity:    0.1, // Start with hammer raised
		forgingFlameVelocity:     0.0,
		forgingFlameTargetToggle: false,
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
			m.projectName = dirName
		} else {
			// Using specific folder - use folder name as default project name
			m.projectName = userPath
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
				Value(&m.setupMode),
		),

		// Group 2: Project name
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Placeholder("my-app").
				Value(&m.projectName).
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
				Value(&m.language),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),

		// Group 4: Framework (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a framework").
				Options(
					huh.NewOption("React", models.FrameworkReact),
					huh.NewOption("Vue", models.FrameworkVue),
					huh.NewOption("Angular", models.FrameworkAngular),
					huh.NewOption("Svelte", models.FrameworkSvelte),
					huh.NewOption("Solid", models.FrameworkSolid),
					huh.NewOption("Vanilla (no framework)", models.FrameworkVanilla),
				).
				Value(&m.framework),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
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
				Value(&m.packageManager),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),

		// Group 6: Styling (only shown in custom mode)
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
				Value(&m.styling),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
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
				Value(&m.uiLibrary),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkReact
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
				Value(&m.uiLibrary),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkVue
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
				Value(&m.uiLibrary),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkAngular
		}),

		// Group 7: Routing for React (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Routing solution").
				Options(
					huh.NewOption("React Router", models.RoutingReactRouter),
					huh.NewOption("TanStack Router", models.RoutingTanStackRouter),
					huh.NewOption("File-based routing", models.RoutingFileBased),
					huh.NewOption("None (single page)", models.RoutingNone),
				).
				Value(&m.routing),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkReact
		}),

		// Group 7b: Routing for Svelte (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Routing solution").
				Options(
					huh.NewOption("SvelteKit (built-in)", models.RoutingSvelteKit),
					huh.NewOption("None (single page)", models.RoutingNone),
				).
				Value(&m.routing),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkSvelte
		}),

		// Group 7c: Routing for Solid (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Routing solution").
				Options(
					huh.NewOption("Solid Router", models.RoutingSolidRouter),
					huh.NewOption("None (single page)", models.RoutingNone),
				).
				Value(&m.routing),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkSolid
		}),

		// Group 8: Testing (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Testing framework").
				Options(
					huh.NewOption("Vitest (fast, Vite-native)", models.TestingVitest),
					huh.NewOption("Jest", models.TestingJest),
					huh.NewOption("None (set up later)", models.TestingNone),
				).
				Value(&m.testing),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),

		// Group 9: State management for React (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("Zustand (lightweight)", models.StateZustand),
					huh.NewOption("Redux Toolkit", models.StateReduxToolkit),
					huh.NewOption("Context API only", models.StateContextAPI),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.stateManagement),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkReact
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
				Value(&m.stateManagement),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkVue
		}),

		// Group 9c: State management for Svelte (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("Svelte Stores (built-in)", models.StateSvelteStores),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.stateManagement),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkSvelte
		}),

		// Group 9d: State management for Solid (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("Solid Stores (built-in)", models.StateSolidStores),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.stateManagement),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkSolid
		}),

		// Group 9e: State management for Angular (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("State management").
				Options(
					huh.NewOption("NgRx", models.StateNgRx),
					huh.NewOption("None", models.StateNone),
				).
				Value(&m.stateManagement),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkAngular
		}),

		// Group 9f: Form Management for React (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Form Management").
				Options(
					huh.NewOption("React Hook Form (recommended)", models.FormReactHookForm),
					huh.NewOption("Formik", models.FormFormik),
					huh.NewOption("TanStack Form", models.FormTanStackForm),
					huh.NewOption("None (native)", models.FormNone),
				).
				Value(&m.formManagement),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkReact
		}),

		// Group 9g: Form Management for Vue (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Form Management").
				Options(
					huh.NewOption("VeeValidate", models.FormVeeValidate),
					huh.NewOption("None (native)", models.FormNone),
				).
				Value(&m.formManagement),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick) || m.framework != models.FrameworkVue
		}),

		// Group 10: Data fetching (only shown in custom mode)
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
				Value(&m.dataFetching),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),

		// Group 10a: Animation (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Animation library").
				Options(
					huh.NewOption("Framer Motion (React)", models.AnimationFramerMotion),
					huh.NewOption("GSAP (any framework)", models.AnimationGSAP),
					huh.NewOption("Auto Animate", models.AnimationAutoAnimate),
					huh.NewOption("React Spring", models.AnimationReactSpring),
					huh.NewOption("Motion One", models.AnimationMotionOne),
					huh.NewOption("None", models.AnimationNone),
				).
				Value(&m.animation),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),

		// Group 10b: Icons (only shown in custom mode)
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
				Value(&m.icons),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),

		// Group 10c: Data Visualization (only shown in custom mode)
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
				Value(&m.dataViz),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),

		// Group 10d: Utilities (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Utility library").
				Options(
					huh.NewOption("date-fns (date manipulation)", models.UtilsDateFns),
					huh.NewOption("Day.js (lightweight dates)", models.UtilsDayJS),
					huh.NewOption("Lodash-es (tree-shakeable)", models.UtilsLodash),
					huh.NewOption("None", models.UtilsNone),
				).
				Value(&m.utilities),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),

		// Group 10e: Internationalization (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Internationalization (i18n)").
				Options(
					huh.NewOption("react-i18next (React)", models.I18nReactI18next),
					huh.NewOption("vue-i18n (Vue)", models.I18nVueI18n),
					huh.NewOption("None", models.I18nNone),
				).
				Value(&m.i18n),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),

		// Group 11: Folder structure (only shown in custom mode)
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Folder structure").
				Options(
					huh.NewOption("Feature-based (recommended)", models.StructureFeatureBased),
					huh.NewOption("Layer-based", models.StructureLayerBased),
				).
				Value(&m.structure),
		).WithHideFunc(func() bool {
			return m.setupMode == string(models.SetupModeQuick)
		}),
		// No confirmation question - review screen handles this
	)

	// Apply custom forge theme to match the rest of the TUI
	form.WithTheme(ForgeTheme())

	return form
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

// applyFormDataToConfig converts form field values to config
func (m *Model) applyFormDataToConfig() {
	m.config.ProjectName = m.projectName

	// If ProjectPath wasn't set (no -path flag), set it based on project name
	if m.config.ProjectPath == "" {
		cwd, err := filepath.Abs(".")
		if err == nil {
			m.config.ProjectPath = filepath.Join(cwd, m.projectName)
		}
	}

	// Apply quick preset if selected
	if m.setupMode == string(models.SetupModeQuick) {
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
	m.config.Language = m.language
	m.config.Framework = m.framework
	m.config.PackageManager = m.packageManager
	m.config.Styling = m.styling
	m.config.UILibrary = m.uiLibrary
	m.config.Testing = m.testing
	m.config.FormManagement = m.formManagement
	m.config.DataFetching = m.dataFetching
	m.config.Animation = m.animation
	m.config.Icons = m.icons
	m.config.DataViz = m.dataViz
	m.config.Utilities = m.utilities
	m.config.I18n = m.i18n
	m.config.Structure = m.structure

	// Handle framework-specific defaults
	switch m.framework {
	case models.FrameworkVanilla:
		m.config.Routing = models.RoutingNone
		m.config.StateManagement = models.StateNone
		m.config.UILibrary = models.UILibraryNone
		m.config.FormManagement = models.FormNone
	case models.FrameworkVue:
		m.config.Routing = models.RoutingVueRouter
		m.config.StateManagement = m.stateManagement
	case models.FrameworkAngular:
		m.config.Routing = models.RoutingAngularRouter
		m.config.StateManagement = m.stateManagement
	case models.FrameworkSvelte:
		m.config.Routing = m.routing
		m.config.StateManagement = m.stateManagement
		m.config.UILibrary = models.UILibraryNone
		m.config.FormManagement = models.FormNone
	case models.FrameworkSolid:
		m.config.Routing = m.routing
		m.config.StateManagement = m.stateManagement
		m.config.UILibrary = models.UILibraryNone
		m.config.FormManagement = models.FormNone
	case models.FrameworkReact:
		m.config.Routing = m.routing
		m.config.StateManagement = m.stateManagement
	}
}
