// Package models defines configuration data structures and constants.
//
// This package contains:
//   - Config struct: Holds all project generation settings
//   - Framework constants: React, Vue, Angular, Svelte, Solid
//   - Language constants: TypeScript, JavaScript
//   - Package manager constants: npm, yarn, pnpm, bun
//   - Styling, routing, testing, and other feature constants
//
// All constants are exported for use in generators and TUI components.
package models

// Config holds all project configuration
type Config struct {
	ProjectName     string
	ProjectPath     string // Absolute path where project will be created
	Language        string
	Framework       string
	PackageManager  string
	Styling         string
	UILibrary       string
	Routing         string
	Testing         string
	StateManagement string
	FormManagement  string
	DataFetching    string
	Animation       string
	Icons           string
	DataViz         string
	Utilities       string
	I18n            string
	Structure       string
	DryRun          bool // Preview mode - show what would be generated without writing files
}

// SetupMode defines quick or custom setup
type SetupMode string

const (
	SetupModeQuick  SetupMode = "quick"
	SetupModeCustom SetupMode = "custom"
)

// Language options
const (
	LangTypeScript = "TypeScript"
	LangJavaScript = "JavaScript"
)

// Framework options
const (
	FrameworkReact   = "React"
	FrameworkVue     = "Vue"
	FrameworkAngular = "Angular"
	FrameworkSvelte  = "Svelte"
	FrameworkSolid   = "Solid"
	FrameworkVanilla = "Vanilla"
)

// Package managers
const (
	PackageManagerNpm  = "npm"
	PackageManagerYarn = "yarn"
	PackageManagerPnpm = "pnpm"
	PackageManagerBun  = "bun"
)

// Styling options
const (
	StylingTailwind   = "Tailwind CSS"
	StylingBootstrap  = "Bootstrap"
	StylingCSSModules = "CSS Modules"
	StylingSass       = "Sass/SCSS"
	StylingStyled     = "Styled Components"
	StylingVanilla    = "Vanilla CSS"
)

// Routing options
const (
	RoutingReactRouter    = "React Router"
	RoutingTanStackRouter = "TanStack Router"
	RoutingFileBased      = "File-based routing"
	RoutingVueRouter      = "Vue Router"
	RoutingAngularRouter  = "Angular Router"
	RoutingSvelteKit      = "SvelteKit"
	RoutingSolidRouter    = "Solid Router"
	RoutingNone           = "None"
)

// Testing options
const (
	TestingVitest = "Vitest"
	TestingJest   = "Jest"
	TestingNone   = "None"
)

// State management options
const (
	StateZustand      = "Zustand"
	StateReduxToolkit = "Redux Toolkit"
	StateContextAPI   = "Context API"
	StatePinia        = "Pinia"
	StateVuex         = "Vuex"
	StateSvelteStores = "Svelte Stores"
	StateSolidStores  = "Solid Stores"
	StateNgRx         = "NgRx"
	StateNone         = "None"
)

// Data fetching options
const (
	DataTanStackQuery = "TanStack Query"
	DataFetchAPI      = "Fetch API"
	DataAxios         = "Axios"
	DataSWR           = "SWR"
	DataNone          = "None"
)

// Project structure options
const (
	StructureFeatureBased = "Feature-based"
	StructureLayerBased   = "Layer-based"
)

// UI Component Library options
const (
	// React
	UILibraryShadcn   = "Shadcn/ui"
	UILibraryMUI      = "Material-UI (MUI)"
	UILibraryChakra   = "Chakra UI"
	UILibraryAntD     = "Ant Design"
	UILibraryHeadless = "Headless UI"
	// Vue
	UILibraryVuetify   = "Vuetify"
	UILibraryPrimeVue  = "PrimeVue"
	UILibraryElementUI = "Element Plus"
	UILibraryNaiveUI   = "Naive UI"
	// Angular
	UILibraryAngularMaterial = "Angular Material"
	UILibraryPrimeNG         = "PrimeNG"
	UILibraryNGZorro         = "NG-ZORRO"
	UILibraryNone            = "None"
)

// Form Management options
const (
	FormReactHookForm = "React Hook Form"
	FormFormik        = "Formik"
	FormTanStackForm  = "TanStack Form"
	FormVeeValidate   = "VeeValidate"
	FormZod           = "Zod (validation)"
	FormYup           = "Yup (validation)"
	FormNone          = "None"
)

// Animation options
const (
	AnimationFramerMotion = "Framer Motion"
	AnimationGSAP         = "GSAP"
	AnimationAutoAnimate  = "Auto Animate"
	AnimationReactSpring  = "React Spring"
	AnimationNone         = "None"
)

// Icon library options
const (
	IconsReactIcons  = "React Icons"
	IconsVueIcons    = "Vue Icons"
	IconsHeroicons   = "Heroicons"
	IconsLucide      = "Lucide"
	IconsFontAwesome = "Font Awesome"
	IconsNone        = "None"
)

// Data Visualization options
const (
	DataVizRecharts = "Recharts"
	DataVizChartJS  = "Chart.js"
	DataVizECharts  = "Apache ECharts"
	DataVizNivo     = "Nivo"
	DataVizNone     = "None"
)

// Utility library options
const (
	UtilsDateFns = "date-fns"
	UtilsDayJS   = "Day.js"
	UtilsLodash  = "Lodash-es"
	UtilsNone    = "None"
)

// Internationalization options
const (
	I18nReactI18next = "react-i18next"
	I18nVueI18n      = "vue-i18n"
	I18nNone         = "None"
)
