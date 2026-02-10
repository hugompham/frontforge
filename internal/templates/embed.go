package templates

import (
	"bytes"
	"embed"
	"fmt"
	"frontforge/internal/models"
	"text/template"
)

//go:embed static/*
var templateFS embed.FS

// TemplateData holds all data for template rendering, including computed fields
type TemplateData struct {
	models.Config
	MountID string // Computed: "app" for Vue, "root" for others
	MainExt string // Computed: main file extension (tsx, ts, js)
	AppExt  string // Computed: app file extension (tsx, vue, svelte, ts, js)
}

// PrepareTemplateData creates TemplateData with computed fields
func PrepareTemplateData(config models.Config) TemplateData {
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
