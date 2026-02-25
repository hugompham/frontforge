package nextjs

import (
	"frontforge/internal/generators/meta"
	"frontforge/internal/generators/shared"
	"frontforge/internal/models"
)

func init() {
	meta.Register(models.FrameworkNextJS, &Generator{})
}

// Generator implements meta.MetaGenerator for Next.js.
type Generator struct{}

func (g *Generator) Scaffold(cfg models.Config) error {
	args := buildScaffoldArgs(cfg)
	return meta.ExecScaffold(models.FrameworkNextJS, cfg.DryRun, "npx", args...)
}

func (g *Generator) PostScaffold(cfg models.Config) error {
	dir := cfg.ProjectPath

	// Merge FrontForge-specific deps
	deps := make(map[string]string)
	devDeps := make(map[string]string)
	scripts := make(map[string]string)

	// ESLint (FrontForge standard)
	devDeps["eslint"] = "^9.39.1"
	devDeps["@eslint/js"] = "^9.39.1"
	devDeps["globals"] = "^15.15.0"
	devDeps["typescript-eslint"] = "^8.56.1"
	devDeps["eslint-plugin-react-hooks"] = "^7.0.1"
	devDeps["eslint-plugin-react-refresh"] = "^0.5.2"
	scripts["lint"] = "eslint ."

	// State management
	switch cfg.StateManagement {
	case models.StateZustand:
		deps["zustand"] = "^5.0.11"
	case models.StateReduxToolkit:
		deps["@reduxjs/toolkit"] = "^2.11.2"
		deps["react-redux"] = "^9.2.0"
	}

	// Data fetching
	switch cfg.DataFetching {
	case models.DataTanStackQuery:
		deps["@tanstack/react-query"] = "^5.90.21"
		devDeps["@tanstack/react-query-devtools"] = "^5.91.3"
	case models.DataAxios:
		deps["axios"] = "^1.13.5"
	case models.DataSWR:
		deps["swr"] = "^2.4.0"
	}

	if len(deps) > 0 || len(devDeps) > 0 || len(scripts) > 0 {
		if err := shared.MergePackageJSON(dir, deps, devDeps, scripts); err != nil {
			return err
		}
	}

	// Testing
	if cfg.Testing == models.TestingVitest {
		if err := shared.ScaffoldVitest(dir, "nextjs"); err != nil {
			return err
		}
	}

	// Feature-based structure
	if cfg.Structure == models.StructureFeatureBased {
		if err := shared.ScaffoldFeatureStructure(dir, "nextjs"); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) SupportedOptions() meta.OptionMatrix {
	return meta.OptionMatrix{
		Styling:         []string{"Tailwind CSS", "CSS Modules", "Sass/SCSS", "Vanilla CSS"},
		Testing:         []string{"Vitest", "Jest", "None"},
		StateManagement: []string{"Zustand", "Redux Toolkit", "Context API", "None"},
		DataFetching:    []string{"TanStack Query", "SWR", "Axios", "Fetch API", "None"},
	}
}

func (g *Generator) ProbeVersion() string {
	return meta.ProbeUpstreamCLI("npx", "create-next-app", "--version")
}

func buildScaffoldArgs(cfg models.Config) []string {
	args := []string{"create-next-app@latest", cfg.ProjectPath}

	// Language
	if cfg.Language == models.LangJavaScript {
		args = append(args, "--js")
	} else {
		args = append(args, "--ts")
	}

	// Tailwind
	if cfg.Styling == models.StylingTailwind {
		args = append(args, "--tailwind")
	} else {
		args = append(args, "--no-tailwind")
	}

	// Standard flags
	args = append(args, "--eslint", "--app", "--src-dir", "--import-alias", "@/*", "--turbopack")

	// Package manager
	switch cfg.PackageManager {
	case models.PackageManagerYarn:
		args = append(args, "--use-yarn")
	case models.PackageManagerPnpm:
		args = append(args, "--use-pnpm")
	case models.PackageManagerBun:
		args = append(args, "--use-bun")
	default:
		args = append(args, "--use-npm")
	}

	// Skip interactive prompts
	args = append(args, "--yes")

	return args
}
