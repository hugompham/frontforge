package generators

import (
	"fmt"
	"frontforge/internal/models"
	"os"
	"path/filepath"
	"strings"
)

// GenerateViteConfig creates vite.config.ts/js
func GenerateViteConfig(config models.Config) string {
	var imports strings.Builder
	var plugins []string

	imports.WriteString("import { defineConfig } from 'vite'\n")

	switch config.Framework {
	case models.FrameworkReact:
		imports.WriteString("import react from '@vitejs/plugin-react'\n")
		plugins = append(plugins, "react()")
	case models.FrameworkVue:
		imports.WriteString("import vue from '@vitejs/plugin-vue'\n")
		plugins = append(plugins, "vue()")
	case models.FrameworkSvelte:
		imports.WriteString("import { svelte } from '@sveltejs/vite-plugin-svelte'\n")
		plugins = append(plugins, "svelte()")
	}

	return fmt.Sprintf(`%s
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [%s],
  server: {
    port: 3000,
    open: true
  }
})
`, imports.String(), strings.Join(plugins, ", "))
}

// TSConfigSet contains all TypeScript configuration files
type TSConfigSet struct {
	Base map[string]interface{}
	App  map[string]interface{}
	Node map[string]interface{}
}

// GenerateTSConfig creates TypeScript configuration
func GenerateTSConfig(config models.Config) TSConfigSet {
	base := map[string]interface{}{
		"files": []string{},
		"references": []map[string]string{
			{"path": "./tsconfig.app.json"},
			{"path": "./tsconfig.node.json"},
		},
	}

	jsx := "preserve"
	if config.Framework == models.FrameworkReact {
		jsx = "react-jsx"
	}

	app := map[string]interface{}{
		"compilerOptions": map[string]interface{}{
			"target":                     "ES2020",
			"useDefineForClassFields":    true,
			"lib":                        []string{"ES2020", "DOM", "DOM.Iterable"},
			"module":                     "ESNext",
			"skipLibCheck":               true,
			"moduleResolution":           "bundler",
			"allowImportingTsExtensions": true,
			"isolatedModules":            true,
			"moduleDetection":            "force",
			"noEmit":                     true,
			"jsx":                        jsx,
			"strict":                     true,
			"noUnusedLocals":             true,
			"noUnusedParameters":         true,
			"noFallthroughCasesInSwitch": true,
		},
		"include": []string{"src"},
	}

	node := map[string]interface{}{
		"compilerOptions": map[string]interface{}{
			"target":                     "ES2022",
			"lib":                        []string{"ES2023"},
			"module":                     "ESNext",
			"skipLibCheck":               true,
			"moduleResolution":           "bundler",
			"allowImportingTsExtensions": true,
			"isolatedModules":            true,
			"moduleDetection":            "force",
			"noEmit":                     true,
			"strict":                     true,
			"noUnusedLocals":             true,
			"noUnusedParameters":         true,
			"noFallthroughCasesInSwitch": true,
		},
		"include": []string{"vite.config.ts"},
	}

	return TSConfigSet{
		Base: base,
		App:  app,
		Node: node,
	}
}

// GenerateProjectStructure creates the directory structure
func GenerateProjectStructure(projectPath string, config models.Config) error {
	// Create base directories
	dirs := []string{
		filepath.Join(projectPath, "src"),
		filepath.Join(projectPath, "public"),
	}

	if config.Structure == models.StructureFeatureBased {
		// Feature-based structure
		dirs = append(dirs,
			filepath.Join(projectPath, "src", "features"),
			filepath.Join(projectPath, "src", "features", "auth"),
			filepath.Join(projectPath, "src", "features", "dashboard"),
			filepath.Join(projectPath, "src", "components"),
			filepath.Join(projectPath, "src", "lib"),
			filepath.Join(projectPath, "src", "hooks"),
		)

		// Create all directories
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		// Create features README
		featureReadme := `# Features

This directory contains feature-based modules. Each feature is self-contained with its own:
- Components
- Hooks
- Services
- Types
- Tests

## Example Structure

` + "```" + `
features/
├── auth/
│   ├── components/
│   ├── hooks/
│   ├── services/
│   └── types/
└── dashboard/
    ├── components/
    ├── hooks/
    └── services/
` + "```" + `
`
		if err := writeFile(filepath.Join(projectPath, "src", "features", "README.md"), featureReadme); err != nil {
			return err
		}
	} else {
		// Layer-based structure
		dirs = append(dirs,
			filepath.Join(projectPath, "src", "components"),
			filepath.Join(projectPath, "src", "pages"),
			filepath.Join(projectPath, "src", "services"),
			filepath.Join(projectPath, "src", "utils"),
			filepath.Join(projectPath, "src", "hooks"),
			filepath.Join(projectPath, "src", "types"),
			filepath.Join(projectPath, "src", "lib"),
		)

		// Create all directories
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
	}

	// Create utils file
	ext := "js"
	if config.Language == models.LangTypeScript {
		ext = "ts"
	}

	typeAnnotation := ""
	if config.Language == models.LangTypeScript {
		typeAnnotation = ": (string | undefined)[]"
	}

	utilsFile := fmt.Sprintf(`export function cn(...classes%s) {
  return classes.filter(Boolean).join(' ');
}
`, typeAnnotation)

	return writeFile(filepath.Join(projectPath, "src", "lib", fmt.Sprintf("utils.%s", ext)), utilsFile)
}
