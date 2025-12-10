package generators

import (
	"frontforge/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SetupProject orchestrates the entire project generation
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

	// Create the project directory if it doesn't exist (for new folder mode)
	// This will do nothing if the directory already exists (current directory mode)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Generate package.json
	packageJSON := GeneratePackageJSON(config)
	if err := writeJSON(filepath.Join(projectPath, "package.json"), packageJSON); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}

	// Generate vite.config for React, Vue, Svelte
	if config.Framework == models.FrameworkReact || config.Framework == models.FrameworkVue || config.Framework == models.FrameworkSvelte {
		viteConfig := GenerateViteConfig(config)
		ext := "js"
		if config.Language == models.LangTypeScript {
			ext = "ts"
		}
		if err := writeFile(filepath.Join(projectPath, fmt.Sprintf("vite.config.%s", ext)), viteConfig); err != nil {
			return fmt.Errorf("failed to write vite.config: %w", err)
		}
	}

	// Generate TypeScript configs
	if config.Language == models.LangTypeScript {
		tsConfigs := GenerateTSConfig(config)
		if err := writeJSON(filepath.Join(projectPath, "tsconfig.json"), tsConfigs.Base); err != nil {
			return fmt.Errorf("failed to write tsconfig.json: %w", err)
		}
		if err := writeJSON(filepath.Join(projectPath, "tsconfig.app.json"), tsConfigs.App); err != nil {
			return fmt.Errorf("failed to write tsconfig.app.json: %w", err)
		}
		if err := writeJSON(filepath.Join(projectPath, "tsconfig.node.json"), tsConfigs.Node); err != nil {
			return fmt.Errorf("failed to write tsconfig.node.json: %w", err)
		}
	}

	// Generate project structure
	if err := GenerateProjectStructure(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate project structure: %w", err)
	}

	// Generate index.html
	indexHTML := GenerateIndexHTML(config)
	if err := writeFile(filepath.Join(projectPath, "index.html"), indexHTML); err != nil {
		return fmt.Errorf("failed to write index.html: %w", err)
	}

	// Generate main entry file
	mainFile := GenerateMainFile(config)
	ext := "jsx"
	if config.Language == models.LangTypeScript {
		ext = "tsx"
	}
	if err := writeFile(filepath.Join(projectPath, "src", fmt.Sprintf("main.%s", ext)), mainFile); err != nil {
		return fmt.Errorf("failed to write main file: %w", err)
	}

	// Generate App component
	appFile := GenerateAppFile(config)
	if err := writeFile(filepath.Join(projectPath, "src", fmt.Sprintf("App.%s", ext)), appFile); err != nil {
		return fmt.Errorf("failed to write App file: %w", err)
	}

	// Generate .gitignore
	gitignore := generateGitignore()
	if err := writeFile(filepath.Join(projectPath, ".gitignore"), gitignore); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	// Generate README
	readme := generateREADME(config)
	if err := writeFile(filepath.Join(projectPath, "README.md"), readme); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// Generate Tailwind config if needed
	if config.Styling == models.StylingTailwind {
		if err := generateTailwindConfig(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate Tailwind config: %w", err)
		}
	}

	// Generate Vitest config if needed
	if config.Testing == models.TestingVitest {
		if err := generateVitestConfig(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate Vitest config: %w", err)
		}
	}

	// Generate ESLint config
	if err := generateESLintConfig(projectPath, config); err != nil {
		return fmt.Errorf("failed to generate ESLint config: %w", err)
	}

	return nil
}

// CleanupCLI removes CLI files after project generation
func CleanupCLI() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	filesToRemove := []string{
		"cli",
		"node_modules",
		"package-lock.json",
		"go.mod",
		"go.sum",
		"internal",
		"main.go",
	}

	for _, file := range filesToRemove {
		filePath := filepath.Join(cwd, file)
		_ = os.RemoveAll(filePath) // Ignore errors
	}

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
` + "```bash\n%s install\n```" + `

### Run development server
` + "```bash\n%s dev\n```" + `

### Build for production
` + "```bash\n%s build\n```" + `

### Run tests
` + "```bash\n%s test\n```" + `

## Project Structure

%s

## Learn More

- [%s Documentation](https://%s.dev)
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
		config.Framework,
	)
}

func generateTailwindConfig(projectPath string, config models.Config) error {
	tailwindConfig := `/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
`
	if err := writeFile(filepath.Join(projectPath, "tailwind.config.js"), tailwindConfig); err != nil {
		return err
	}

	postcssConfig := `export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
`
	if err := writeFile(filepath.Join(projectPath, "postcss.config.js"), postcssConfig); err != nil {
		return err
	}

	indexCSS := `@tailwind base;
@tailwind components;
@tailwind utilities;
`
	return writeFile(filepath.Join(projectPath, "src", "index.css"), indexCSS)
}

func generateVitestConfig(projectPath string, config models.Config) error {
	ext := "js"
	if config.Language == models.LangTypeScript {
		ext = "ts"
	}

	vitestConfig := fmt.Sprintf(`import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.%s',
  },
});
`, ext)

	if err := writeFile(filepath.Join(projectPath, fmt.Sprintf("vitest.config.%s", ext)), vitestConfig); err != nil {
		return err
	}

	// Create test directory
	testDir := filepath.Join(projectPath, "src", "test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return err
	}

	setupFile := fmt.Sprintf(`import { expect, afterEach } from 'vitest';
import { cleanup } from '@testing-library/react';
import * as matchers from '@testing-library/jest-dom/matchers';

expect.extend(matchers);

afterEach(() => {
  cleanup();
});
`)
	return writeFile(filepath.Join(testDir, fmt.Sprintf("setup.%s", ext)), setupFile)
}

func generateESLintConfig(projectPath string, config models.Config) error {
	eslintConfig := `import js from '@eslint/js'
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
	return writeFile(filepath.Join(projectPath, "eslint.config.js"), eslintConfig)
}
