// Package generators handles project file and structure generation.
//
// This package orchestrates the creation of frontend projects by generating
// all necessary files and directories based on the user's configuration.
//
// Key features:
//   - Automatic cleanup on generation failure (rollback pattern)
//   - Path validation and safety checks
//   - Support for React, Vue, Angular, Svelte, and Solid frameworks
//   - TypeScript and JavaScript support
//   - Vite-based build configuration
//
// Main entry point is SetupProject(), which coordinates all file generation
// and ensures atomic operations with cleanup on error.
package generators

import (
	"encoding/json"
	"fmt"
	"frontforge/internal/models"
	"frontforge/internal/templates"
	"os"
	"path/filepath"
)

// SetupProject orchestrates the entire project generation
// On error, automatically cleans up any partially created files
// If config.DryRun is true, prints a manifest without writing files
func SetupProject(config models.Config) error {
	// Use the project path from config (set by CLI flags)
	projectPath := config.ProjectPath
	if projectPath == "" {
		// Fallback to current directory if not set (shouldn't happen with new flag system)
		var err error
		projectPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Dry-run mode: collect manifest instead of writing files
	var manifest *DryRunManifest
	if config.DryRun {
		manifest = NewDryRunManifest(projectPath, config.ProjectName)
		defer manifest.Print()
	}

	// Track all created paths for cleanup on failure
	var createdPaths []string
	success := false

	// Cleanup function that removes created files in reverse order
	cleanup := func() {
		for i := len(createdPaths) - 1; i >= 0; i-- {
			_ = os.RemoveAll(createdPaths[i])
		}
	}

	// Defer cleanup if generation fails
	defer func() {
		if !success && !config.DryRun {
			cleanup()
		}
	}()

	// Create the project directory if it doesn't exist (for new folder mode)
	// This will do nothing if the directory already exists (current directory mode)
	dirExisted := false
	if !config.DryRun {
		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			if err := os.MkdirAll(projectPath, 0755); err != nil {
				return fmt.Errorf("failed to create project directory: %w", err)
			}
			createdPaths = append(createdPaths, projectPath)
		} else {
			dirExisted = true
		}
	}

	// Helper to track file creation
	trackPath := func(path string) {
		// Only track if directory didn't exist before (avoid deleting user's directory)
		if !dirExisted && !config.DryRun {
			createdPaths = append(createdPaths, path)
		}
	}

	// Helper to write or collect files
	writeOrCollect := func(path, content string) error {
		if manifest != nil {
			manifest.AddFile(path, content)
			return nil
		}
		return writeFile(path, content)
	}

	// Helper to write or collect JSON
	writeOrCollectJSON := func(path string, data interface{}) error {
		bytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		content := string(bytes)
		if manifest != nil {
			manifest.AddFile(path, content)
			return nil
		}
		return writeJSON(path, data)
	}

	// Helper to create or collect directories
	mkdirOrCollect := func(path string) error {
		if manifest != nil {
			manifest.AddDir(path)
			return nil
		}
		return os.MkdirAll(path, 0755)
	}

	// Generate package.json
	packageJSON := GeneratePackageJSON(config)
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if err := writeOrCollectJSON(packageJSONPath, packageJSON); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}
	trackPath(packageJSONPath)

	// Generate vite.config for all Vite-based frameworks
	if config.Framework == models.FrameworkReact || config.Framework == models.FrameworkVue || config.Framework == models.FrameworkSvelte || config.Framework == models.FrameworkSolid || config.Framework == models.FrameworkAngular || config.Framework == models.FrameworkVanilla {
		viteConfig := GenerateViteConfig(config)
		ext := "js"
		if config.Language == models.LangTypeScript {
			ext = "ts"
		}
		if err := writeOrCollect(filepath.Join(projectPath, fmt.Sprintf("vite.config.%s", ext)), viteConfig); err != nil {
			return fmt.Errorf("failed to write vite.config: %w", err)
		}
	}

	// Generate TypeScript configs
	if config.Language == models.LangTypeScript {
		tsConfigs := GenerateTSConfig(config)
		if err := writeOrCollectJSON(filepath.Join(projectPath, "tsconfig.json"), tsConfigs.Base); err != nil {
			return fmt.Errorf("failed to write tsconfig.json: %w", err)
		}
		if err := writeOrCollectJSON(filepath.Join(projectPath, "tsconfig.app.json"), tsConfigs.App); err != nil {
			return fmt.Errorf("failed to write tsconfig.app.json: %w", err)
		}
		if err := writeOrCollectJSON(filepath.Join(projectPath, "tsconfig.node.json"), tsConfigs.Node); err != nil {
			return fmt.Errorf("failed to write tsconfig.node.json: %w", err)
		}
	}

	// Generate project structure
	if err := GenerateProjectStructure(projectPath, config, mkdirOrCollect, writeOrCollect); err != nil {
		return fmt.Errorf("failed to generate project structure: %w", err)
	}

	// Generate index.html
	indexHTML := GenerateIndexHTML(config)
	if err := writeOrCollect(filepath.Join(projectPath, "index.html"), indexHTML); err != nil {
		return fmt.Errorf("failed to write index.html: %w", err)
	}

	// Generate vite.svg favicon
	viteSVG, err := templates.RenderStatic("static/vite.svg")
	if err != nil {
		return fmt.Errorf("failed to generate vite.svg: %w", err)
	}
	if err := writeOrCollect(filepath.Join(projectPath, "public", "vite.svg"), viteSVG); err != nil {
		return fmt.Errorf("failed to write vite.svg: %w", err)
	}

	// Generate main entry file
	mainFile := GenerateMainFile(config)
	ext := "jsx"
	if config.Language == models.LangTypeScript {
		ext = "tsx"
	}

	// Angular uses .ts extension always (no .tsx)
	if config.Framework == models.FrameworkAngular {
		ext = "ts"
	}

	if err := writeOrCollect(filepath.Join(projectPath, "src", fmt.Sprintf("main.%s", ext)), mainFile); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	// Generate App component
	appFile := GenerateAppFile(config)

	// Angular components go in app/ directory
	if config.Framework == models.FrameworkAngular {
		// Create app directory
		if err := os.MkdirAll(filepath.Join(projectPath, "src", "app"), 0755); err != nil {
			return fmt.Errorf("failed to create app directory: %w", err)
		}
		if err := writeOrCollect(filepath.Join(projectPath, "src", "app", "app.component.ts"), appFile); err != nil {
			return fmt.Errorf("failed to write App component: %w", err)
		}
	} else {
		if err := writeOrCollect(filepath.Join(projectPath, "src", fmt.Sprintf("App.%s", ext)), appFile); err != nil {
			return fmt.Errorf("failed to write App file: %w", err)
		}
	}

	// Generate .gitignore
	gitignore, err := templates.RenderStatic("static/gitignore.tmpl")
	if err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}
	if err := writeOrCollect(filepath.Join(projectPath, ".gitignore"), gitignore); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	// Generate README
	readme, err := templates.Render("static/README.md.tmpl", config)
	if err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}
	if err := writeOrCollect(filepath.Join(projectPath, "README.md"), readme); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// Generate Tailwind config if needed
	if config.Styling == models.StylingTailwind {
		indexCSS, err := templates.RenderStatic("static/index.css")
		if err != nil {
			return fmt.Errorf("failed to read Tailwind CSS import file: %w", err)
		}
		if err := writeOrCollect(filepath.Join(projectPath, "src", "index.css"), indexCSS); err != nil {
			return fmt.Errorf("failed to write index.css: %w", err)
		}
	}

	// Generate CSS Modules example if needed
	if config.Styling == models.StylingCSSModules {
		appModuleCSS, err := templates.RenderStatic("static/App.module.css")
		if err != nil {
			return fmt.Errorf("failed to generate CSS Modules example: %w", err)
		}
		if err := writeOrCollect(filepath.Join(projectPath, "src", "App.module.css"), appModuleCSS); err != nil {
			return fmt.Errorf("failed to write CSS Modules example: %w", err)
		}
	}

	// Generate Sass example if needed
	if config.Styling == models.StylingSass {
		stylesScss, err := templates.RenderStatic("static/styles.scss")
		if err != nil {
			return fmt.Errorf("failed to generate Sass example: %w", err)
		}
		if err := writeOrCollect(filepath.Join(projectPath, "src", "styles.scss"), stylesScss); err != nil {
			return fmt.Errorf("failed to write Sass example: %w", err)
		}
	}

	// Generate Vitest config if needed
	if config.Testing == models.TestingVitest {
		ext := "js"
		if config.Language == models.LangTypeScript {
			ext = "ts"
		}

		// Generate vitest config
		vitestConfig, err := templates.RenderVitestConfig(config)
		if err != nil {
			return fmt.Errorf("failed to generate Vitest config: %w", err)
		}
		if err := writeOrCollect(filepath.Join(projectPath, fmt.Sprintf("vitest.config.%s", ext)), vitestConfig); err != nil {
			return fmt.Errorf("failed to write Vitest config: %w", err)
		}

		// Create test directory
		testDir := filepath.Join(projectPath, "src", "test")
		if err := mkdirOrCollect(testDir); err != nil {
			return fmt.Errorf("failed to create test directory: %w", err)
		}

		// Generate test setup
		setupFile, err := templates.RenderVitestSetup(config)
		if err != nil {
			return fmt.Errorf("failed to generate Vitest setup: %w", err)
		}
		if err := writeOrCollect(filepath.Join(testDir, fmt.Sprintf("setup.%s", ext)), setupFile); err != nil {
			return fmt.Errorf("failed to write Vitest setup: %w", err)
		}
	}

	// Generate ESLint config
	eslintConfig, err := templates.RenderESLintConfig(config)
	if err != nil {
		return fmt.Errorf("failed to generate ESLint config: %w", err)
	}
	if err := writeOrCollect(filepath.Join(projectPath, "eslint.config.js"), eslintConfig); err != nil {
		return fmt.Errorf("failed to write ESLint config: %w", err)
	}

	// Run post-generation validation
	validationResults := ValidateProject(projectPath, config)
	failedChecks := 0
	for _, result := range validationResults {
		if !result.Passed {
			failedChecks++
			if !config.DryRun {
				fmt.Printf("  Warning: %s - %s\n", result.Check, result.Message)
			}
		}
	}

	if failedChecks > 0 && !config.DryRun {
		fmt.Printf("\nValidation completed with %d warning(s). Project may not work correctly.\n", failedChecks)
	}

	// Mark generation as successful - prevents cleanup
	success = true
	return nil
}

// Helper functions
func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func writeJSON(path string, data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}
