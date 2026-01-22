package tui

import "frontforge/internal/models"

// ============================================================================
// FORGING THEME ICONS - ASCII-based symbols for terminal compatibility
// ============================================================================

// ForgeIcons contains all forge-themed ASCII symbols
var ForgeIcons = struct {
	// Core Forging Tools
	Hammer    string
	Anvil     string
	Furnace   string
	Blueprint string
	Tongs     string
	Sparks    string
	Fire      string

	// Materials & States
	IronIngot string
	Steel     string
	Metal     string
	Molten    string

	// Process Indicators
	Heating   string
	Forging   string
	Tempering string
	Finished  string

	// Status Indicators
	Success string
	Error   string
	Warning string
	Info    string

	// Progress Indicators
	InProgress string
	Pending    string
	Complete   string
}{
	// Core Forging Tools
	Hammer:    "[H]",
	Anvil:     "[A]",
	Furnace:   "[F]",
	Blueprint: "[B]",
	Tongs:     "[T]",
	Sparks:    "[*]",
	Fire:      "[~]",

	// Materials & States
	IronIngot: "[=]",
	Steel:     "[-]",
	Metal:     "[#]",
	Molten:    "[o]",

	// Process Indicators
	Heating:   "[^]",
	Forging:   "[!]",
	Tempering: "[v]",
	Finished:  "[+]",

	// Status Indicators (ASCII-only for maximum compatibility)
	Success: "[+]",
	Error:   "[X]",
	Warning: "[!]",
	Info:    "[i]",

	// Progress Indicators
	InProgress: "[>]",
	Pending:    "[ ]",
	Complete:   "[+]",
}

// ============================================================================
// FRAMEWORK ICONS - ASCII representations of frameworks as forging tools
// ============================================================================

// Each framework is represented by a tool in the forge metaphor:
// - React: Hammer (building blocks, component-based)
// - Vue: Chisel (sculpting, reactive)
// - Angular: Anvil (structured, framework-heavy)
// - Svelte: Lightning (fast, compiled)
// - Solid: Diamond (solid, crystallized)
// - Vanilla: Flame (pure, raw)

var frameworkForgeIcons = map[string]string{
	models.FrameworkReact:   "[H]", // Hammer - building blocks
	models.FrameworkVue:     "[C]", // Chisel - sculpting
	models.FrameworkAngular: "[A]", // Anvil - structured
	models.FrameworkSvelte:  "[L]", // Lightning - fast
	models.FrameworkSolid:   "[D]", // Diamond - solid
	models.FrameworkVanilla: "[F]", // Flame - pure
}

// ============================================================================
// PACKAGE MANAGER ICONS - Material types in the forge
// ============================================================================

// Each package manager is represented by a material type:
// - npm: Standard Ore (common, widely used)
// - yarn: Thread/Weave (interconnected)
// - pnpm: Precision Metal (efficient, exact)
// - bun: Bread/Fuel (fast, energy)

var packageManagerForgeIcons = map[string]string{
	models.PackageManagerNpm:  "[N]", // Npm - standard ore
	models.PackageManagerYarn: "[Y]", // Yarn - threaded
	models.PackageManagerPnpm: "[P]", // Pnpm - precision
	models.PackageManagerBun:  "[B]", // Bun - fuel/energy
}

// ============================================================================
// STYLING ICONS - Finishing techniques in the forge
// ============================================================================

// Each styling solution is a finishing technique:
// - Tailwind: Polish (utility-based refinement)
// - Bootstrap: Template (pre-made patterns)
// - Sass: Advanced tools (enhanced capabilities)
// - Styled Components: Custom forge (component-specific)
// - CSS Modules: Modular pieces (scoped)
// - Vanilla: Raw finish (pure CSS)

var stylingForgeIcons = map[string]string{
	models.StylingTailwind:   "[T]", // Tailwind - polish
	models.StylingBootstrap:  "[B]", // Bootstrap - template
	models.StylingSass:       "[S]", // Sass - advanced
	models.StylingStyled:     "[C]", // Styled - custom
	models.StylingCSSModules: "[M]", // Modules - modular
	models.StylingVanilla:    "[V]", // Vanilla - raw
}

// ============================================================================
// UI LIBRARY ICONS - Pre-forged components
// ============================================================================

var uiLibraryForgeIcons = map[string]string{
	// React Libraries
	models.UILibraryShadcn:   "[S]", // Shadcn
	models.UILibraryMUI:      "[M]", // Material UI
	models.UILibraryChakra:   "[C]", // Chakra
	models.UILibraryAntD:     "[A]", // Ant Design
	models.UILibraryHeadless: "[H]", // Headless UI

	// Vue Libraries
	models.UILibraryVuetify:   "[V]", // Vuetify
	models.UILibraryPrimeVue:  "[P]", // PrimeVue
	models.UILibraryElementUI: "[E]", // Element Plus
	models.UILibraryNaiveUI:   "[N]", // Naive UI

	// Angular Libraries
	models.UILibraryAngularMaterial: "[M]", // Angular Material
	models.UILibraryPrimeNG:         "[P]", // PrimeNG
	models.UILibraryNGZorro:         "[Z]", // NG-ZORRO
}

// ============================================================================
// TESTING ICONS - Quality control in the forge
// ============================================================================

var testingForgeIcons = map[string]string{
	models.TestingVitest: "[V]", // Vitest - modern testing
	models.TestingJest:   "[J]", // Jest - traditional testing
}

// ============================================================================
// STATE MANAGEMENT ICONS - Blueprint organization
// ============================================================================

var stateForgeIcons = map[string]string{
	models.StateZustand:      "[Z]", // Zustand - lightweight
	models.StateReduxToolkit: "[R]", // Redux - structured
	models.StatePinia:        "[P]", // Pinia - Vue state
	models.StateVuex:         "[V]", // Vuex - legacy Vue
	models.StateNgRx:         "[N]", // NgRx - Angular state
	models.StateSvelteStores: "[S]", // Svelte stores
	models.StateSolidStores:  "[S]", // Solid stores
}

// ============================================================================
// ICON GETTER FUNCTIONS
// ============================================================================

// GetFrameworkForgeIcon returns the forging-themed icon for a framework
func GetFrameworkForgeIcon(framework string) string {
	return "" // Icons removed for cleaner UI
}

// GetPackageManagerForgeIcon returns the forging-themed icon for a package manager
func GetPackageManagerForgeIcon(pm string) string {
	return "" // Icons removed for cleaner UI
}

// GetStylingForgeIcon returns the forging-themed icon for a styling solution
func GetStylingForgeIcon(styling string) string {
	return "" // Icons removed for cleaner UI
}

// GetUILibraryForgeIcon returns the forging-themed icon for a UI library
func GetUILibraryForgeIcon(lib string) string {
	return "" // Icons removed for cleaner UI
}

// GetTestingForgeIcon returns the forging-themed icon for a testing framework
func GetTestingForgeIcon(testing string) string {
	return "" // Icons removed for cleaner UI
}

// GetStateForgeIcon returns the forging-themed icon for state management
func GetStateForgeIcon(state string) string {
	return "" // Icons removed for cleaner UI
}

// ============================================================================
// SPECIAL FORGE SYMBOLS - For visual elements in the TUI
// ============================================================================

// GetForgeTitle returns an ASCII art title for "FRONTFORGE"
func GetForgeTitle() string {
	return `F R O N T F O R G E`
}

// GetForgeWelcomeArt returns ASCII art for the welcome screen
func GetForgeWelcomeArt() string {
	return `
          THE FORGE

       ~~~  Molten Hot  ~~~

          Raw -> Forged
`
}
