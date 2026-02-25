package meta

import "frontforge/internal/models"

// OptionMatrix describes which options a meta-framework supports.
// nil fields are hidden entirely in the TUI.
type OptionMatrix struct {
	Styling         []string
	Testing         []string
	StateManagement []string // nil = hidden
	DataFetching    []string // nil = hidden
}

// MetaGenerator defines the interface for meta-framework generators.
type MetaGenerator interface {
	// Scaffold runs the upstream CLI non-interactively.
	Scaffold(cfg models.Config) error

	// PostScaffold applies FrontForge additions (testing, state, etc.).
	PostScaffold(cfg models.Config) error

	// SupportedOptions returns what TUI should show for this framework.
	SupportedOptions() OptionMatrix

	// ProbeVersion detects installed upstream CLI version; returns "" if not found.
	ProbeVersion() string
}

// generators maps framework constants to their MetaGenerator implementations.
// Populated by init() in each per-framework package.
var generators = map[string]MetaGenerator{}

// Register adds a MetaGenerator for a framework. Called from per-framework init().
func Register(framework string, gen MetaGenerator) {
	generators[framework] = gen
}

// Get returns the MetaGenerator for a framework, or nil if not found.
func Get(framework string) (MetaGenerator, bool) {
	g, ok := generators[framework]
	return g, ok
}

// RunMetaScaffold is the high-level entry point called from SetupProject.
// It runs Scaffold + PostScaffold for a meta-framework config.
func RunMetaScaffold(cfg models.Config) error {
	gen, ok := Get(cfg.Framework)
	if !ok {
		return &ScaffoldError{
			Framework: cfg.Framework,
			Command:   "",
			ExitCode:  -1,
			Stderr:    "no generator registered for framework",
		}
	}

	if !cfg.NoScaffold {
		if err := gen.Scaffold(cfg); err != nil {
			return err
		}
	}

	// In dry-run mode, scaffold is skipped so there's no project directory
	// to apply post-scaffold transforms to. Just report what would happen.
	if cfg.DryRun {
		return nil
	}

	return gen.PostScaffold(cfg)
}
