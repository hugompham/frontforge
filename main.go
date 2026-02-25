package main

import (
	"flag"
	"fmt"
	"frontforge/internal/generators"
	"frontforge/internal/models"
	"frontforge/internal/preflight"
	"frontforge/internal/tui"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Define command-line flags
	var projectPath string
	var showHelp bool
	var quickMode bool
	var dryRun bool
	var autoInstall bool
	var projectName string
	var framework string
	var language string
	var packageManager string
	var styling string

	flag.StringVar(&projectPath, "path", "", "Project path (use '.' for current directory, or specify a folder name)")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.BoolVar(&showHelp, "h", false, "Show help information (shorthand)")

	// Non-interactive flags
	flag.BoolVar(&quickMode, "quick", false, "Use quick preset (React + TypeScript + Tailwind) and skip interactive mode")
	flag.BoolVar(&dryRun, "dry-run", false, "Preview mode: show what files would be generated without writing them")
	flag.BoolVar(&autoInstall, "install", false, "Automatically run package manager install after generation")
	flag.StringVar(&projectName, "name", "", "Project name (required for non-interactive mode)")
	flag.StringVar(&framework, "framework", "", "Framework: react, vue, angular, svelte, solid, vanilla, nextjs, astro, sveltekit")
	flag.StringVar(&language, "lang", "", "Language: ts, js")
	flag.StringVar(&packageManager, "pm", "", "Package manager: npm, yarn, pnpm, bun")
	flag.StringVar(&styling, "styling", "", "Styling: tailwind, bootstrap, css-modules, sass, styled, vanilla")

	// Additional option flags
	var testing string
	var stateManagement string
	var dataFetching string
	flag.StringVar(&testing, "testing", "", "Testing: vitest, jest, playwright, none")
	flag.StringVar(&stateManagement, "state", "", "State management: zustand, redux, pinia, svelte-stores, context, none")
	flag.StringVar(&dataFetching, "data", "", "Data fetching: tanstack-query, swr, axios, fetch, none")

	// Meta-framework debugging
	var noScaffold bool
	flag.BoolVar(&noScaffold, "no-scaffold", false, "Skip upstream CLI scaffold (meta-frameworks only, for debugging)")

	flag.Parse()

	// Show help if requested
	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// Check if running in non-interactive mode
	if quickMode || projectName != "" {
		runNonInteractive(projectPath, projectName, quickMode, dryRun, autoInstall, noScaffold, framework, language, packageManager, styling, testing, stateManagement, dataFetching)
		return
	}

	// Resolve the absolute project path
	var absPath string
	var userPath string
	var err error

	if projectPath == "" {
		// No path specified - will create a new folder based on project name
		// Path will be resolved later after user enters project name in TUI
		absPath = ""
		userPath = ""
	} else if projectPath == "." {
		// Use current directory
		absPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		userPath = "."
	} else {
		// Create new directory with specified path
		if filepath.IsAbs(projectPath) {
			absPath = filepath.Clean(projectPath)
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				os.Exit(1)
			}
			absPath = filepath.Join(cwd, projectPath)
		}
		userPath = projectPath
	}

	// Create the Bubbletea program with project path
	p := tea.NewProgram(tui.NewModelWithPath(absPath, userPath))

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

// runNonInteractive generates a project without the interactive TUI
func runNonInteractive(projectPath, projectName string, quickMode, dryRun, autoInstall, noScaffold bool, framework, language, packageManager, styling, testing, stateManagement, dataFetching string) {
	// Validate project name is provided
	if projectName == "" {
		fmt.Println("Error: -name flag is required for non-interactive mode")
		fmt.Println("Usage: frontforge -quick -name my-project")
		os.Exit(1)
	}

	// Validate project name format
	if !isValidProjectName(projectName) {
		fmt.Println("Error: Invalid project name. Use only letters, numbers, hyphens, and underscores.")
		os.Exit(1)
	}

	// Start with quick preset as base
	config := models.QuickPreset()
	config.ProjectName = projectName
	config.DryRun = dryRun
	config.AutoInstall = autoInstall
	config.NoScaffold = noScaffold

	// Apply overrides if provided
	if framework != "" {
		if fw := parseFramework(framework); fw != "" {
			config.Framework = fw
			// Adjust framework-specific defaults
			adjustFrameworkDefaults(&config)
		} else {
			fmt.Printf("Error: Invalid framework '%s'. Valid options: react, vue, angular, svelte, solid, vanilla, nextjs, astro, sveltekit\n", framework)
			os.Exit(1)
		}
	}

	if language != "" {
		if lang := parseLanguage(language); lang != "" {
			config.Language = lang
		} else {
			fmt.Printf("Error: Invalid language '%s'. Valid options: ts, js\n", language)
			os.Exit(1)
		}
	}

	if packageManager != "" {
		if pm := parsePackageManager(packageManager); pm != "" {
			config.PackageManager = pm
		} else {
			fmt.Printf("Error: Invalid package manager '%s'. Valid options: npm, yarn, pnpm, bun\n", packageManager)
			os.Exit(1)
		}
	}

	if styling != "" {
		if st := parseStyling(styling); st != "" {
			config.Styling = st
		} else {
			fmt.Printf("Error: Invalid styling '%s'. Valid options: tailwind, bootstrap, css-modules, sass, styled, vanilla\n", styling)
			os.Exit(1)
		}
	}

	if testing != "" {
		if t := parseTesting(testing); t != "" {
			config.Testing = t
		} else {
			fmt.Printf("Error: Invalid testing '%s'. Valid options: vitest, jest, playwright, none\n", testing)
			os.Exit(1)
		}
	}

	if stateManagement != "" {
		if sm := parseStateManagement(stateManagement); sm != "" {
			config.StateManagement = sm
		} else {
			fmt.Printf("Error: Invalid state management '%s'. Valid options: zustand, redux, pinia, svelte-stores, context, none\n", stateManagement)
			os.Exit(1)
		}
	}

	if dataFetching != "" {
		if df := parseDataFetching(dataFetching); df != "" {
			config.DataFetching = df
		} else {
			fmt.Printf("Error: Invalid data fetching '%s'. Valid options: tanstack-query, swr, axios, fetch, none\n", dataFetching)
			os.Exit(1)
		}
	}

	// Resolve project path
	var absPath string
	var err error

	if projectPath == "" {
		// Create folder with project name
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		absPath = filepath.Join(cwd, projectName)
	} else if projectPath == "." {
		absPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
	} else {
		if filepath.IsAbs(projectPath) {
			absPath = filepath.Clean(projectPath)
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				os.Exit(1)
			}
			absPath = filepath.Join(cwd, projectPath)
		}
	}

	config.ProjectPath = absPath

	// Validate framework + library compatibility
	if err := validateCompatibility(&config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Print configuration summary
	fmt.Println()
	fmt.Println("FrontForge - Non-Interactive Mode")
	fmt.Println()
	fmt.Printf("  Project:   %s\n", config.ProjectName)
	fmt.Printf("  Path:      %s\n", config.ProjectPath)
	fmt.Printf("  Framework: %s\n", config.Framework)
	fmt.Printf("  Language:  %s\n", config.Language)
	fmt.Printf("  Styling:   %s\n", config.Styling)
	fmt.Printf("  Package:   %s\n", config.PackageManager)
	fmt.Println()

	// Run preflight checks
	fmt.Println("Running pre-flight checks...")
	results := preflight.RunAllChecks(config)

	for _, check := range results.Checks {
		if check.Passed {
			fmt.Printf("  [OK] %s\n", check.Name)
		} else {
			fmt.Printf("  [FAIL] %s: %s\n", check.Name, check.Message)
			if check.Suggestion != "" {
				fmt.Printf("    → %s\n", check.Suggestion)
			}
		}
	}

	if results.FatalError {
		fmt.Println()
		fmt.Println("Pre-flight checks failed. Please resolve the issues above.")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Generating project...")

	// Generate the project
	if err := generators.SetupProject(config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Run install if requested (only for actual generation, not dry-run)
	if config.AutoInstall && !config.DryRun {
		fmt.Println()
		fmt.Printf("Running %s install...\n", config.PackageManager)
		fmt.Println()

		if err := generators.RunInstall(config.ProjectPath, config); err != nil {
			fmt.Println()
			fmt.Printf("Warning: Install failed: %v\n", err)
			fmt.Println("You can run the install manually with:")
			fmt.Printf("  cd %s\n", config.ProjectName)
			fmt.Printf("  %s install\n", config.PackageManager)
		} else {
			fmt.Println()
			fmt.Println("Dependencies installed successfully!")
		}
	}

	// Success message
	fmt.Println()
	fmt.Println("Project created successfully!")
	fmt.Println()
	fmt.Println("Next steps:")

	// Only show cd command if project was created in a subdirectory
	cwd, _ := filepath.Abs(".")
	if config.ProjectPath != cwd {
		fmt.Printf("  cd %s\n", config.ProjectName)
	}

	// Skip install step if auto-install was used successfully
	if !config.AutoInstall || config.DryRun {
		fmt.Printf("  %s install\n", config.PackageManager)
	}

	fmt.Printf("  %s run dev\n", getRunCommand(config.PackageManager))
	fmt.Println()
}

// isValidProjectName checks if the project name contains only valid characters
func isValidProjectName(name string) bool {
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return len(name) > 0
}

// validateCompatibility checks for incompatible framework + library combinations
func validateCompatibility(config *models.Config) error {
	// Meta-frameworks have built-in routing; skip most compatibility checks
	if models.IsMetaFramework(config.Framework) {
		return nil
	}

	// Validate routing compatibility
	switch config.Framework {
	case models.FrameworkReact:
		if config.Routing != models.RoutingReactRouter && config.Routing != models.RoutingTanStackRouter && config.Routing != models.RoutingFileBased && config.Routing != models.RoutingNone {
			return fmt.Errorf("routing '%s' is not compatible with React", config.Routing)
		}
	case models.FrameworkVue:
		if config.Routing != models.RoutingVueRouter && config.Routing != models.RoutingNone {
			return fmt.Errorf("routing '%s' is not compatible with Vue", config.Routing)
		}
	case models.FrameworkSvelte:
		if config.Routing != models.RoutingSvelteKit && config.Routing != models.RoutingNone {
			return fmt.Errorf("routing '%s' is not compatible with Svelte", config.Routing)
		}
	case models.FrameworkSolid:
		if config.Routing != models.RoutingSolidRouter && config.Routing != models.RoutingNone {
			return fmt.Errorf("routing '%s' is not compatible with Solid", config.Routing)
		}
	}

	// Validate state management compatibility
	switch config.Framework {
	case models.FrameworkReact:
		if config.StateManagement != models.StateZustand && config.StateManagement != models.StateReduxToolkit && config.StateManagement != models.StateContextAPI && config.StateManagement != models.StateNone {
			return fmt.Errorf("state management '%s' is not compatible with React", config.StateManagement)
		}
	case models.FrameworkVue:
		if config.StateManagement != models.StatePinia && config.StateManagement != models.StateVuex && config.StateManagement != models.StateNone {
			return fmt.Errorf("state management '%s' is not compatible with Vue", config.StateManagement)
		}
	case models.FrameworkSvelte:
		if config.StateManagement != models.StateSvelteStores && config.StateManagement != models.StateNone {
			return fmt.Errorf("state management '%s' is not compatible with Svelte", config.StateManagement)
		}
	case models.FrameworkSolid:
		if config.StateManagement != models.StateSolidStores && config.StateManagement != models.StateNone {
			return fmt.Errorf("state management '%s' is not compatible with Solid", config.StateManagement)
		}
	}

	// Validate UI library compatibility
	reactUILibs := []string{models.UILibraryMUI, models.UILibraryChakra, models.UILibraryAntD, models.UILibraryShadcn, models.UILibraryHeadless}
	vueUILibs := []string{models.UILibraryVuetify, models.UILibraryPrimeVue, models.UILibraryElementUI, models.UILibraryNaiveUI}
	angularUILibs := []string{models.UILibraryAngularMaterial, models.UILibraryPrimeNG, models.UILibraryNGZorro}

	if config.Framework != models.FrameworkReact {
		for _, lib := range reactUILibs {
			if config.UILibrary == lib {
				return fmt.Errorf("UI library '%s' is only compatible with React", lib)
			}
		}
	}

	if config.Framework != models.FrameworkVue {
		for _, lib := range vueUILibs {
			if config.UILibrary == lib {
				return fmt.Errorf("UI library '%s' is only compatible with Vue", lib)
			}
		}
	}

	if config.Framework != models.FrameworkAngular {
		for _, lib := range angularUILibs {
			if config.UILibrary == lib {
				return fmt.Errorf("UI library '%s' is only compatible with Angular", lib)
			}
		}
	}

	// Validate form management compatibility
	reactFormLibs := []string{models.FormReactHookForm, models.FormFormik, models.FormTanStackForm}
	vueFormLibs := []string{models.FormVeeValidate}

	if config.Framework != models.FrameworkReact {
		for _, lib := range reactFormLibs {
			if config.FormManagement == lib {
				return fmt.Errorf("form management '%s' is only compatible with React", lib)
			}
		}
	}

	if config.Framework != models.FrameworkVue {
		for _, lib := range vueFormLibs {
			if config.FormManagement == lib {
				return fmt.Errorf("form management '%s' is only compatible with Vue", lib)
			}
		}
	}

	return nil
}

// parseFramework converts a short framework name to the full constant
func parseFramework(input string) string {
	switch strings.ToLower(input) {
	case "react":
		return models.FrameworkReact
	case "vue":
		return models.FrameworkVue
	case "angular":
		return models.FrameworkAngular
	case "svelte":
		return models.FrameworkSvelte
	case "solid":
		return models.FrameworkSolid
	case "vanilla":
		return models.FrameworkVanilla
	case "nextjs", "next", "next.js":
		return models.FrameworkNextJS
	case "astro":
		return models.FrameworkAstro
	case "sveltekit", "svelte-kit":
		return models.FrameworkSvelteKit
	default:
		return ""
	}
}

// parseLanguage converts a short language name to the full constant
func parseLanguage(input string) string {
	switch strings.ToLower(input) {
	case "ts", "typescript":
		return models.LangTypeScript
	case "js", "javascript":
		return models.LangJavaScript
	default:
		return ""
	}
}

// parsePackageManager converts a package manager name to the constant
func parsePackageManager(input string) string {
	switch strings.ToLower(input) {
	case "npm":
		return models.PackageManagerNpm
	case "yarn":
		return models.PackageManagerYarn
	case "pnpm":
		return models.PackageManagerPnpm
	case "bun":
		return models.PackageManagerBun
	default:
		return ""
	}
}

// parseStyling converts a styling name to the full constant
func parseStyling(input string) string {
	switch strings.ToLower(input) {
	case "tailwind", "tailwindcss":
		return models.StylingTailwind
	case "bootstrap":
		return models.StylingBootstrap
	case "css-modules", "cssmodules":
		return models.StylingCSSModules
	case "sass", "scss":
		return models.StylingSass
	case "styled", "styled-components":
		return models.StylingStyled
	case "vanilla", "css":
		return models.StylingVanilla
	default:
		return ""
	}
}

// parseTesting converts a testing name to the full constant
func parseTesting(input string) string {
	switch strings.ToLower(input) {
	case "vitest":
		return models.TestingVitest
	case "jest":
		return models.TestingJest
	case "playwright":
		return models.TestingPlaywright
	case "none":
		return models.TestingNone
	default:
		return ""
	}
}

// parseStateManagement converts a state management name to the full constant
func parseStateManagement(input string) string {
	switch strings.ToLower(input) {
	case "zustand":
		return models.StateZustand
	case "redux", "redux-toolkit":
		return models.StateReduxToolkit
	case "pinia":
		return models.StatePinia
	case "svelte-stores", "svelte":
		return models.StateSvelteStores
	case "context", "context-api":
		return models.StateContextAPI
	case "none":
		return models.StateNone
	default:
		return ""
	}
}

// parseDataFetching converts a data fetching name to the full constant
func parseDataFetching(input string) string {
	switch strings.ToLower(input) {
	case "tanstack-query", "tanstack", "react-query":
		return models.DataTanStackQuery
	case "swr":
		return models.DataSWR
	case "axios":
		return models.DataAxios
	case "fetch", "fetch-api":
		return models.DataFetchAPI
	case "none":
		return models.DataNone
	default:
		return ""
	}
}

// adjustFrameworkDefaults sets sensible defaults for non-React frameworks
func adjustFrameworkDefaults(config *models.Config) {
	switch config.Framework {
	case models.FrameworkVue:
		config.Routing = models.RoutingVueRouter
		config.StateManagement = models.StatePinia
		config.UILibrary = models.UILibraryVuetify
		config.FormManagement = models.FormVeeValidate
		config.DataFetching = models.DataAxios
		config.Icons = models.IconsVueIcons
		config.I18n = models.I18nVueI18n
		config.Animation = models.AnimationAutoAnimate
	case models.FrameworkAngular:
		config.Routing = models.RoutingAngularRouter
		config.StateManagement = models.StateNgRx
		config.UILibrary = models.UILibraryAngularMaterial
		config.FormManagement = models.FormNone
		config.DataFetching = models.DataFetchAPI
		config.Animation = models.AnimationNone
		config.Icons = models.IconsNone
		config.I18n = models.I18nNone
	case models.FrameworkSvelte:
		config.Routing = models.RoutingSvelteKit
		config.StateManagement = models.StateSvelteStores
		config.UILibrary = models.UILibraryNone
		config.FormManagement = models.FormNone
		config.DataFetching = models.DataFetchAPI
		config.Animation = models.AnimationAutoAnimate
		config.Icons = models.IconsLucide
		config.I18n = models.I18nNone
	case models.FrameworkSolid:
		config.Routing = models.RoutingSolidRouter
		config.StateManagement = models.StateSolidStores
		config.UILibrary = models.UILibraryNone
		config.FormManagement = models.FormNone
		config.DataFetching = models.DataFetchAPI
		config.Animation = models.AnimationFramerMotion
		config.Icons = models.IconsLucide
		config.I18n = models.I18nNone
	case models.FrameworkVanilla:
		config.Routing = models.RoutingNone
		config.StateManagement = models.StateNone
		config.UILibrary = models.UILibraryNone
		config.FormManagement = models.FormNone
		config.DataFetching = models.DataFetchAPI
		config.Animation = models.AnimationNone
		config.Icons = models.IconsNone
		config.I18n = models.I18nNone
	case models.FrameworkNextJS:
		config.Routing = models.RoutingNextJSAppRouter
		config.StateManagement = models.StateNone
		config.UILibrary = models.UILibraryNone
		config.FormManagement = models.FormNone
		config.DataFetching = models.DataFetchAPI
		config.Animation = models.AnimationNone
		config.Icons = models.IconsLucide
		config.I18n = models.I18nNone
	case models.FrameworkAstro:
		config.Routing = models.RoutingAstroPages
		config.StateManagement = models.StateNone
		config.UILibrary = models.UILibraryNone
		config.FormManagement = models.FormNone
		config.DataFetching = models.DataFetchAPI
		config.Animation = models.AnimationNone
		config.Icons = models.IconsNone
		config.I18n = models.I18nNone
	case models.FrameworkSvelteKit:
		config.Routing = models.RoutingSvelteKit
		config.StateManagement = models.StateSvelteStores
		config.UILibrary = models.UILibraryNone
		config.FormManagement = models.FormNone
		config.DataFetching = models.DataFetchAPI
		config.Animation = models.AnimationNone
		config.Icons = models.IconsLucide
		config.I18n = models.I18nNone
	}
}

// getRunCommand returns the appropriate run command for the package manager
func getRunCommand(pm string) string {
	if pm == models.PackageManagerNpm {
		return "npm"
	}
	return pm
}

// printHelp displays comprehensive help information
func printHelp() {
	fmt.Println("FRONTFORGE - Modern Frontend Project Scaffolding")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  frontforge [options]")
	fmt.Println()
	fmt.Println("INTERACTIVE MODE (default):")
	fmt.Println("  Run without flags to use the interactive TUI for full configuration.")
	fmt.Println()
	fmt.Println("NON-INTERACTIVE MODE:")
	fmt.Println("  Use -quick and -name flags for instant project generation.")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  -h, -help        Show this help information")
	fmt.Println()
	fmt.Println("  Path:")
	fmt.Println("    -path <path>   Project path (optional)")
	fmt.Println("                   If not specified, creates folder with project name")
	fmt.Println("                   Use '.' for current directory")
	fmt.Println()
	fmt.Println("  Non-Interactive:")
	fmt.Println("    -quick         Use quick preset and skip interactive mode")
	fmt.Println("    -name <name>   Project name (required for non-interactive)")
	fmt.Println("    -framework     Framework: react, vue, angular, svelte, solid, vanilla,")
	fmt.Println("                             nextjs, astro, sveltekit")
	fmt.Println("    -lang          Language: ts, js (default: ts)")
	fmt.Println("    -pm            Package manager: npm, yarn, pnpm, bun (default: npm)")
	fmt.Println("    -styling       Styling: tailwind, bootstrap, css-modules, sass, styled, vanilla")
	fmt.Println("    -testing       Testing: vitest, jest, playwright, none")
	fmt.Println("    -state         State: zustand, redux, pinia, svelte-stores, context, none")
	fmt.Println("    -data          Data fetching: tanstack-query, swr, axios, fetch, none")
	fmt.Println("    -no-scaffold   Skip upstream CLI (meta-frameworks only, for debugging)")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  Interactive mode:")
	fmt.Println("    frontforge")
	fmt.Println("    frontforge -path .")
	fmt.Println()
	fmt.Println("  Quick generation (React + TypeScript + Tailwind):")
	fmt.Println("    frontforge -quick -name my-app")
	fmt.Println()
	fmt.Println("  Vue project with pnpm:")
	fmt.Println("    frontforge -quick -name my-vue-app -framework vue -pm pnpm")
	fmt.Println()
	fmt.Println("  Svelte project with JavaScript:")
	fmt.Println("    frontforge -quick -name my-svelte-app -framework svelte -lang js")
	fmt.Println()
	fmt.Println("  Next.js project:")
	fmt.Println("    frontforge -quick -name my-next-app -framework nextjs")
	fmt.Println()
	fmt.Println("  Astro project:")
	fmt.Println("    frontforge -quick -name my-astro-app -framework astro")
	fmt.Println()
	fmt.Println("  SvelteKit project:")
	fmt.Println("    frontforge -quick -name my-sveltekit-app -framework sveltekit")
	fmt.Println()
	fmt.Println("  Project in current directory:")
	fmt.Println("    frontforge -quick -name my-app -path .")
	fmt.Println()
	fmt.Println("QUICK PRESET INCLUDES:")
	fmt.Println("  Framework:        React (with Vite)")
	fmt.Println("  Language:         TypeScript")
	fmt.Println("  Styling:          Tailwind CSS")
	fmt.Println("  UI Library:       Shadcn/ui")
	fmt.Println("  Routing:          React Router")
	fmt.Println("  State:            Zustand")
	fmt.Println("  Data Fetching:    TanStack Query")
	fmt.Println("  Forms:            React Hook Form")
	fmt.Println("  Testing:          Vitest")
	fmt.Println("  Animation:        Framer Motion")
	fmt.Println("  Icons:            Heroicons")
	fmt.Println("  Utilities:        date-fns")
	fmt.Println()
}
