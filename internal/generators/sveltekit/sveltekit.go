package sveltekit

import (
	"frontforge/internal/generators/meta"
	"frontforge/internal/generators/shared"
	"frontforge/internal/models"
)

func init() {
	meta.Register(models.FrameworkSvelteKit, &Generator{})
}

// Generator implements meta.MetaGenerator for SvelteKit.
type Generator struct{}

func (g *Generator) Scaffold(cfg models.Config) error {
	// Step 1: Create project with sv create
	args := buildCreateArgs(cfg)
	if err := meta.ExecScaffold(models.FrameworkSvelteKit, cfg.DryRun, "npx", args...); err != nil {
		return err
	}

	// Step 2: Add add-ons via sv add
	addOns := buildAddOns(cfg)
	if len(addOns) > 0 {
		svArgs := []string{"sv", "add"}
		svArgs = append(svArgs, addOns...)
		if err := meta.ExecInDir(cfg.ProjectPath, models.FrameworkSvelteKit, cfg.DryRun, "npx", svArgs...); err != nil {
			return err
		}
	}

	return nil
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

	// State management
	if cfg.StateManagement == models.StateSvelteStores {
		// Built-in, no extra deps needed
	}

	// Data fetching
	switch cfg.DataFetching {
	case models.DataTanStackQuery:
		deps["@tanstack/svelte-query"] = "^6.0.18"
	case models.DataAxios:
		deps["axios"] = "^1.13.5"
	}

	if len(deps) > 0 || len(devDeps) > 0 || len(scripts) > 0 {
		if err := shared.MergePackageJSON(dir, deps, devDeps, scripts); err != nil {
			return err
		}
	}

	// Feature-based structure
	if cfg.Structure == models.StructureFeatureBased {
		if err := shared.ScaffoldFeatureStructure(dir, "sveltekit"); err != nil {
			return err
		}
	}

	// Install dependencies
	installCmd := cfg.PackageManager
	if installCmd == "" {
		installCmd = "npm"
	}
	return meta.ExecInDir(dir, models.FrameworkSvelteKit, cfg.DryRun, installCmd, "install")
}

func (g *Generator) SupportedOptions() meta.OptionMatrix {
	return meta.OptionMatrix{
		Styling:         []string{"Tailwind CSS", "CSS Modules", "Sass/SCSS", "Vanilla CSS"},
		Testing:         []string{"Vitest", "Playwright", "None"},
		StateManagement: []string{"Svelte Stores", "None"},
		DataFetching:    []string{"TanStack Query", "Fetch API", "None"},
	}
}

func (g *Generator) ProbeVersion() string {
	return meta.ProbeUpstreamCLI("npx", "sv", "--version")
}

func buildCreateArgs(cfg models.Config) []string {
	args := []string{"sv", "create", cfg.ProjectPath, "--template", "minimal"}

	if cfg.Language == models.LangJavaScript {
		args = append(args, "--no-types")
	} else {
		args = append(args, "--types", "ts")
	}

	args = append(args, "--no-add-ons")

	return args
}

func buildAddOns(cfg models.Config) []string {
	var addOns []string

	if cfg.Styling == models.StylingTailwind {
		addOns = append(addOns, "tailwindcss")
	}

	switch cfg.Testing {
	case models.TestingVitest:
		addOns = append(addOns, "vitest")
	case models.TestingPlaywright:
		addOns = append(addOns, "playwright")
	}

	// Always add eslint and prettier
	addOns = append(addOns, "eslint", "prettier")

	return addOns
}
