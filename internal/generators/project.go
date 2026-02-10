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
	viteSVG := generateViteSVG()
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
	gitignore := generateGitignore()
	if err := writeOrCollect(filepath.Join(projectPath, ".gitignore"), gitignore); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	// Generate README
	readme := generateREADME(config)
	if err := writeOrCollect(filepath.Join(projectPath, "README.md"), readme); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// Generate Tailwind config if needed
	if config.Styling == models.StylingTailwind {
		if err := generateTailwindConfig(projectPath, config, writeOrCollect); err != nil {
			return fmt.Errorf("failed to generate Tailwind config: %w", err)
		}
	}

	// Generate CSS Modules example if needed
	if config.Styling == models.StylingCSSModules {
		if err := generateCSSModulesExample(projectPath, config, writeOrCollect); err != nil {
			return fmt.Errorf("failed to generate CSS Modules example: %w", err)
		}
	}

	// Generate Sass example if needed
	if config.Styling == models.StylingSass {
		if err := generateSassExample(projectPath, config, writeOrCollect); err != nil {
			return fmt.Errorf("failed to generate Sass example: %w", err)
		}
	}

	// Generate Vitest config if needed
	if config.Testing == models.TestingVitest {
		if err := generateVitestConfig(projectPath, config, mkdirOrCollect, writeOrCollect); err != nil {
			return fmt.Errorf("failed to generate Vitest config: %w", err)
		}
	}

	// Generate ESLint config
	if err := generateESLintConfig(projectPath, config, writeOrCollect); err != nil {
		return fmt.Errorf("failed to generate ESLint config: %w", err)
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

func generateGitignore() string {
	return `# Logs
logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*
lerna-debug.log*

node_modules
dist
dist-ssr
*.local

# Editor directories and files
.vscode/*
!.vscode/extensions.json
.idea
.DS_Store
*.suo
*.ntvs*
*.njsproj
*.sln
*.sw?
`
}

func generateREADME(config models.Config) string {
	pmRun := config.PackageManager
	if config.PackageManager == models.PackageManagerNpm {
		pmRun = "npm run"
	}

	structureExample := ""
	if config.Structure == models.StructureFeatureBased {
		ext := "jsx"
		if config.Language == models.LangTypeScript {
			ext = "tsx"
		}
		structureExample = fmt.Sprintf("```\nsrc/\n├── features/          # Feature-based modules\n│   ├── auth/\n│   ├── dashboard/\n│   └── users/\n├── components/        # Shared components\n├── lib/              # Utilities and helpers\n├── App.%s\n└── main.%s\n```", ext, ext)
	} else {
		ext := "jsx"
		if config.Language == models.LangTypeScript {
			ext = "tsx"
		}
		structureExample = fmt.Sprintf("```\nsrc/\n├── components/       # UI components\n├── pages/           # Page components\n├── services/        # API services\n├── utils/           # Utility functions\n├── App.%s\n└── main.%s\n```", ext, ext)
	}

	return fmt.Sprintf(`# %s

This project was created with frontforge.

## Tech Stack

- **Language**: %s
- **Framework**: %s
- **Build Tool**: Vite
- **Styling**: %s
- **Routing**: %s
- **Testing**: %s
- **State Management**: %s
- **Data Fetching**: %s

## Getting Started

### Install dependencies
`+"```bash\n%s install\n```"+`

### Run development server
`+"```bash\n%s dev\n```"+`

### Build for production
`+"```bash\n%s build\n```"+`

### Run tests
`+"```bash\n%s test\n```"+`

## Project Structure

%s

## Learn More

- [%s Documentation](%s)
- [Vite Documentation](https://vitejs.dev)
`,
		config.ProjectName,
		config.Language,
		config.Framework,
		config.Styling,
		config.Routing,
		config.Testing,
		config.StateManagement,
		config.DataFetching,
		config.PackageManager,
		pmRun,
		pmRun,
		pmRun,
		structureExample,
		config.Framework,
		getFrameworkDocURL(config.Framework),
	)
}

// getFrameworkDocURL returns the correct documentation URL for each framework
func getFrameworkDocURL(framework string) string {
	switch framework {
	case models.FrameworkReact:
		return "https://react.dev"
	case models.FrameworkVue:
		return "https://vuejs.org"
	case models.FrameworkAngular:
		return "https://angular.dev"
	case models.FrameworkSvelte:
		return "https://svelte.dev"
	case models.FrameworkSolid:
		return "https://www.solidjs.com"
	case models.FrameworkVanilla:
		return "https://developer.mozilla.org/en-US/docs/Web/JavaScript"
	default:
		return "https://vitejs.dev"
	}
}

func generateTailwindConfig(
	projectPath string,
	config models.Config,
	writeFunc func(string, string) error,
) error {
	// Tailwind CSS 4 uses CSS-first configuration via @import
	// No tailwind.config.js or postcss.config.js needed
	indexCSS := `@import "tailwindcss";
`
	return writeFunc(filepath.Join(projectPath, "src", "index.css"), indexCSS)
}

func generateCSSModulesExample(
	projectPath string,
	config models.Config,
	writeFunc func(string, string) error,
) error {
	// Generate App.module.css with example styles
	appModuleCSS := `.app {
  text-align: center;
  padding: 2rem;
}

.title {
  font-size: 2rem;
  font-weight: bold;
  margin-bottom: 1rem;
  color: #333;
}

.button {
  padding: 0.5rem 1rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.button:hover {
  background-color: #0056b3;
}
`
	return writeFunc(filepath.Join(projectPath, "src", "App.module.css"), appModuleCSS)
}

func generateSassExample(
	projectPath string,
	config models.Config,
	writeFunc func(string, string) error,
) error {
	// Generate styles.scss with example styles and Sass features
	stylesScss := `// Variables
$primary-color: #007bff;
$primary-hover: #0056b3;
$text-color: #333;
$spacing: 1rem;

// Mixins
@mixin button-styles {
  padding: $spacing * 0.5 $spacing;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

// Styles
.app {
  text-align: center;
  padding: $spacing * 2;

  .title {
    font-size: 2rem;
    font-weight: bold;
    margin-bottom: $spacing;
    color: $text-color;
  }

  .button {
    @include button-styles;
    background-color: $primary-color;
    color: white;

    &:hover {
      background-color: $primary-hover;
    }
  }
}
`
	return writeFunc(filepath.Join(projectPath, "src", "styles.scss"), stylesScss)
}

func generateVitestConfig(
	projectPath string,
	config models.Config,
	mkdirFunc func(string) error,
	writeFunc func(string, string) error,
) error {
	ext := "js"
	if config.Language == models.LangTypeScript {
		ext = "ts"
	}

	// Generate framework-specific plugin import and usage
	var pluginImport, pluginUsage string
	switch config.Framework {
	case models.FrameworkReact:
		pluginImport = "import react from '@vitejs/plugin-react';"
		pluginUsage = "react()"
	case models.FrameworkVue:
		pluginImport = "import vue from '@vitejs/plugin-vue';"
		pluginUsage = "vue()"
	case models.FrameworkSvelte:
		pluginImport = "import { svelte } from '@sveltejs/vite-plugin-svelte';"
		pluginUsage = "svelte()"
	default:
		pluginImport = ""
		pluginUsage = ""
	}

	vitestConfig := fmt.Sprintf(`import { defineConfig } from 'vitest/config';
%s

export default defineConfig({
  plugins: [%s],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.%s',
  },
});
`, pluginImport, pluginUsage, ext)

	if err := writeFunc(filepath.Join(projectPath, fmt.Sprintf("vitest.config.%s", ext)), vitestConfig); err != nil {
		return err
	}

	// Create test directory
	testDir := filepath.Join(projectPath, "src", "test")
	if err := mkdirFunc(testDir); err != nil {
		return err
	}

	// Generate framework-specific test setup
	var setupFile string
	switch config.Framework {
	case models.FrameworkReact:
		setupFile = `import { expect, afterEach } from 'vitest';
import { cleanup } from '@testing-library/react';
import * as matchers from '@testing-library/jest-dom/matchers';

expect.extend(matchers);

afterEach(() => {
  cleanup();
});
`
	case models.FrameworkVue:
		setupFile = `import { expect } from 'vitest';
import * as matchers from '@testing-library/jest-dom/matchers';

expect.extend(matchers);
`
	case models.FrameworkSvelte:
		setupFile = `import { expect, afterEach } from 'vitest';
import { cleanup } from '@testing-library/svelte';
import * as matchers from '@testing-library/jest-dom/matchers';

expect.extend(matchers);

afterEach(() => {
  cleanup();
});
`
	default:
		// Vanilla JS/TS - minimal setup
		setupFile = `import { expect } from 'vitest';
import * as matchers from '@testing-library/jest-dom/matchers';

expect.extend(matchers);
`
	}

	return writeFunc(filepath.Join(testDir, fmt.Sprintf("setup.%s", ext)), setupFile)
}

func generateESLintConfig(
	projectPath string,
	config models.Config,
	writeFunc func(string, string) error,
) error {
	var eslintConfig string

	switch config.Framework {
	case models.FrameworkReact:
		eslintConfig = `import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'

export default tseslint.config(
  { ignores: ['dist'] },
  {
    extends: [js.configs.recommended, ...tseslint.configs.recommended],
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
    plugins: {
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
    },
    rules: {
      ...reactHooks.configs.recommended.rules,
      'react-refresh/only-export-components': [
        'warn',
        { allowConstantExport: true },
      ],
    },
  },
)
`
	case models.FrameworkVue:
		eslintConfig = `import js from '@eslint/js'
import globals from 'globals'
import pluginVue from 'eslint-plugin-vue'
import tseslint from 'typescript-eslint'

export default tseslint.config(
  { ignores: ['dist'] },
  {
    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommended,
      ...pluginVue.configs['flat/recommended'],
    ],
    files: ['**/*.{ts,tsx,vue}'],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
  },
)
`
	default:
		// Svelte, Solid, Vanilla - basic TypeScript + ESLint
		eslintConfig = `import js from '@eslint/js'
import globals from 'globals'
import tseslint from 'typescript-eslint'

export default tseslint.config(
  { ignores: ['dist'] },
  {
    extends: [js.configs.recommended, ...tseslint.configs.recommended],
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
  },
)
`
	}

	return writeFunc(filepath.Join(projectPath, "eslint.config.js"), eslintConfig)
}

func generateViteSVG() string {
	// Simple Vite-themed SVG logo (lightning bolt in circle)
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <circle cx="12" cy="12" r="10" stroke="#646cff" fill="none"/>
  <path d="M13 2L3 14h8l-1 8 10-12h-8l1-8z" fill="#646cff"/>
</svg>
`
}
