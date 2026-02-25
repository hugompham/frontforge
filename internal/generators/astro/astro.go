package astro

import (
	"fmt"
	"frontforge/internal/generators/meta"
	"frontforge/internal/generators/shared"
	"frontforge/internal/models"
	"os"
	"path/filepath"
)

func init() {
	meta.Register(models.FrameworkAstro, &Generator{})
}

// Generator implements meta.MetaGenerator for Astro.
type Generator struct{}

func (g *Generator) Scaffold(cfg models.Config) error {
	args := buildScaffoldArgs(cfg)
	return meta.ExecScaffold(models.FrameworkAstro, cfg.DryRun, "npm", args...)
}

func (g *Generator) PostScaffold(cfg models.Config) error {
	dir := cfg.ProjectPath

	deps := make(map[string]string)
	devDeps := make(map[string]string)
	scripts := make(map[string]string)

	// ESLint
	devDeps["eslint"] = "^9.39.1"
	devDeps["@eslint/js"] = "^9.39.1"
	devDeps["globals"] = "^15.15.0"
	devDeps["typescript-eslint"] = "^8.56.1"
	scripts["lint"] = "eslint ."

	// Tailwind (Astro 5.2+: use @tailwindcss/vite, not @astrojs/tailwind)
	if cfg.Styling == models.StylingTailwind {
		devDeps["tailwindcss"] = "^4.2.1"
		devDeps["@tailwindcss/vite"] = "^4.2.1"

		// Create global CSS with Tailwind import
		stylesDir := filepath.Join(dir, "src", "styles")
		if err := os.MkdirAll(stylesDir, 0755); err != nil {
			return fmt.Errorf("failed to create styles directory: %w", err)
		}
		globalCSS := "@import \"tailwindcss\";\n"
		if err := os.WriteFile(filepath.Join(stylesDir, "global.css"), []byte(globalCSS), 0644); err != nil {
			return fmt.Errorf("failed to write global.css: %w", err)
		}

		// Add @tailwindcss/vite to astro.config
		astroConfig := `import { defineConfig } from 'astro/config'
import tailwindcss from '@tailwindcss/vite'

// https://astro.build/config
export default defineConfig({
  vite: {
    plugins: [tailwindcss()],
  },
})
`
		if err := os.WriteFile(filepath.Join(dir, "astro.config.mjs"), []byte(astroConfig), 0644); err != nil {
			return fmt.Errorf("failed to write astro.config.mjs: %w", err)
		}
	}

	if len(deps) > 0 || len(devDeps) > 0 || len(scripts) > 0 {
		if err := shared.MergePackageJSON(dir, deps, devDeps, scripts); err != nil {
			return err
		}
	}

	// Testing
	if cfg.Testing == models.TestingVitest {
		if err := shared.ScaffoldVitest(dir, "astro"); err != nil {
			return err
		}
	}

	// Feature-based structure
	if cfg.Structure == models.StructureFeatureBased {
		if err := shared.ScaffoldFeatureStructure(dir, "astro"); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) SupportedOptions() meta.OptionMatrix {
	return meta.OptionMatrix{
		Styling:         []string{"Tailwind CSS", "CSS Modules", "Sass/SCSS", "Vanilla CSS"},
		Testing:         []string{"Vitest", "None"},
		StateManagement: nil, // Hidden in TUI
		DataFetching:    nil, // Hidden in TUI
	}
}

func (g *Generator) ProbeVersion() string {
	return meta.ProbeUpstreamCLI("npm", "create", "astro@latest", "--", "--version")
}

func buildScaffoldArgs(cfg models.Config) []string {
	args := []string{"create", "astro@latest", cfg.ProjectPath, "--"}

	args = append(args, "--template", "minimal")

	if cfg.Language == models.LangJavaScript {
		args = append(args, "--typescript", "relaxed")
	} else {
		args = append(args, "--typescript", "strict")
	}

	args = append(args, "--install", "--git", "--skip-houston")

	return args
}
