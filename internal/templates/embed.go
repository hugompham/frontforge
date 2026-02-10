package templates

import (
	"bytes"
	"embed"
	"fmt"
	"frontforge/internal/models"
	"text/template"
)

//go:embed static/* config/*
var templateFS embed.FS

// TemplateData holds all data for template rendering, including computed fields
type TemplateData struct {
	models.Config
	MountID          string // Computed: "app" for Vue, "root" for others
	MainExt          string // Computed: main file extension (tsx, ts, js)
	AppExt           string // Computed: app file extension (tsx, vue, svelte, ts, js)
	PmRun            string // Computed: package manager run command
	StructureExample string // Computed: project structure markdown
	FrameworkDocURL  string // Computed: framework documentation URL
	VitestExt        string // Computed: vitest config/setup extension (ts or js)
}

// PrepareTemplateData creates TemplateData with computed fields
func PrepareTemplateData(config models.Config) TemplateData {
	// Validate ProjectName to prevent template injection
	// Only allow alphanumeric, hyphens, and underscores
	for _, r := range config.ProjectName {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_') {
			config.ProjectName = "my-app" // safe fallback
			break
		}
	}

	data := TemplateData{
		Config: config,
	}

	// Compute MountID
	if config.Framework == models.FrameworkVue {
		data.MountID = "app"
	} else {
		data.MountID = "root"
	}

	// Compute MainExt
	data.MainExt = computeMainFileExtension(config)

	// Compute AppExt
	data.AppExt = computeAppFileExtension(config)

	// Compute PmRun
	// yarn/pnpm/bun support bare commands (e.g., "yarn dev"),
	// but npm requires "npm run dev"
	data.PmRun = config.PackageManager
	if config.PackageManager == models.PackageManagerNpm {
		data.PmRun = "npm run"
	}

	// Compute StructureExample
	data.StructureExample = computeStructureExample(config)

	// Compute FrameworkDocURL
	data.FrameworkDocURL = computeFrameworkDocURL(config.Framework)

	// Compute VitestExt
	if config.Language == models.LangTypeScript {
		data.VitestExt = "ts"
	} else {
		data.VitestExt = "js"
	}

	return data
}

// Render renders a template with the given config data
func Render(templatePath string, config models.Config) (string, error) {
	// Read the template file
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Prepare template data with computed fields
	data := PrepareTemplateData(config)

	// Parse and execute the template
	tmpl, err := template.New(templatePath).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// computeMainFileExtension returns the file extension for main file
func computeMainFileExtension(config models.Config) string {
	isTS := config.Language == models.LangTypeScript

	// Angular always uses .ts (no .tsx)
	if config.Framework == models.FrameworkAngular {
		return "ts"
	}

	// Vanilla uses plain .js or .ts (no JSX)
	if config.Framework == models.FrameworkVanilla {
		if isTS {
			return "ts"
		}
		return "js"
	}

	// React and Solid use .tsx for TypeScript
	if config.Framework == models.FrameworkReact || config.Framework == models.FrameworkSolid {
		if isTS {
			return "tsx"
		}
		return "jsx"
	}

	// Vue and Svelte use .ts or .js
	if isTS {
		return "ts"
	}
	return "js"
}

// computeAppFileExtension returns the file extension for App file
func computeAppFileExtension(config models.Config) string {
	switch config.Framework {
	case models.FrameworkVue:
		return "vue"
	case models.FrameworkSvelte:
		return "svelte"
	case models.FrameworkAngular:
		return "ts"
	default:
		if config.Language == models.LangTypeScript {
			return "tsx"
		}
		return "jsx"
	}
}

// RenderStatic returns static file content without template processing
func RenderStatic(filePath string) (string, error) {
	content, err := templateFS.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read static file %s: %w", filePath, err)
	}
	return string(content), nil
}

// RenderESLintConfig renders the appropriate ESLint config based on framework
func RenderESLintConfig(config models.Config) (string, error) {
	var templatePath string

	switch config.Framework {
	case models.FrameworkReact:
		templatePath = "config/eslint-react.tmpl"
	case models.FrameworkVue:
		templatePath = "config/eslint-vue.tmpl"
	default:
		// Svelte, Solid, Vanilla, Angular - use default config
		templatePath = "config/eslint-default.tmpl"
	}

	return RenderStatic(templatePath)
}

// RenderVitestConfig renders the appropriate Vitest config based on framework
func RenderVitestConfig(config models.Config) (string, error) {
	var templatePath string

	switch config.Framework {
	case models.FrameworkReact:
		templatePath = "config/vitest-react.tmpl"
	case models.FrameworkVue:
		templatePath = "config/vitest-vue.tmpl"
	case models.FrameworkSvelte:
		templatePath = "config/vitest-svelte.tmpl"
	default:
		templatePath = "config/vitest-default.tmpl"
	}

	return Render(templatePath, config)
}

// RenderVitestSetup renders the appropriate Vitest setup file based on framework.
// Uses RenderStatic because setup files contain no template directives.
func RenderVitestSetup(config models.Config) (string, error) {
	var templatePath string

	switch config.Framework {
	case models.FrameworkReact:
		templatePath = "config/vitest-setup-react.tmpl"
	case models.FrameworkVue:
		templatePath = "config/vitest-setup-vue.tmpl"
	case models.FrameworkSvelte:
		templatePath = "config/vitest-setup-svelte.tmpl"
	default:
		templatePath = "config/vitest-setup-default.tmpl"
	}

	return RenderStatic(templatePath)
}

// computeStructureExample returns the project structure markdown for README
func computeStructureExample(config models.Config) string {
	ext := "jsx"
	if config.Language == models.LangTypeScript {
		ext = "tsx"
	}

	if config.Structure == models.StructureFeatureBased {
		return fmt.Sprintf("```\nsrc/\n├── features/          # Feature-based modules\n│   ├── auth/\n│   ├── dashboard/\n│   └── users/\n├── components/        # Shared components\n├── lib/              # Utilities and helpers\n├── App.%s\n└── main.%s\n```", ext, ext)
	}

	return fmt.Sprintf("```\nsrc/\n├── components/       # UI components\n├── pages/           # Page components\n├── services/        # API services\n├── utils/           # Utility functions\n├── App.%s\n└── main.%s\n```", ext, ext)
}

// computeFrameworkDocURL returns the documentation URL for a framework
func computeFrameworkDocURL(framework string) string {
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
