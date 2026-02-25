package sveltekit

import (
	"frontforge/internal/generators/meta"
	"frontforge/internal/models"
	"testing"
)

// Compile-time interface compliance check.
var _ meta.MetaGenerator = (*Generator)(nil)

func TestBuildCreateArgs(t *testing.T) {
	tests := []struct {
		name     string
		cfg      models.Config
		wantArgs []string
	}{
		{
			name: "TypeScript project",
			cfg: models.Config{
				ProjectPath: "/tmp/svelte-ts",
				Language:    models.LangTypeScript,
			},
			wantArgs: []string{
				"sv", "create", "/tmp/svelte-ts",
				"--template", "minimal",
				"--types", "ts",
				"--no-add-ons",
			},
		},
		{
			name: "JavaScript project",
			cfg: models.Config{
				ProjectPath: "/tmp/svelte-js",
				Language:    models.LangJavaScript,
			},
			wantArgs: []string{
				"sv", "create", "/tmp/svelte-js",
				"--template", "minimal",
				"--no-types",
				"--no-add-ons",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCreateArgs(tt.cfg)
			assertArgsEqual(t, got, tt.wantArgs)
		})
	}
}

func TestBuildAddOns(t *testing.T) {
	tests := []struct {
		name    string
		cfg     models.Config
		wantOns []string
	}{
		{
			name: "Tailwind and Vitest",
			cfg: models.Config{
				Styling: models.StylingTailwind,
				Testing: models.TestingVitest,
			},
			wantOns: []string{"tailwindcss", "vitest", "eslint", "prettier"},
		},
		{
			name: "no Tailwind, no Vitest",
			cfg: models.Config{
				Styling: models.StylingVanilla,
				Testing: models.TestingNone,
			},
			wantOns: []string{"eslint", "prettier"},
		},
		{
			name: "Tailwind only, no Vitest",
			cfg: models.Config{
				Styling: models.StylingTailwind,
				Testing: models.TestingNone,
			},
			wantOns: []string{"tailwindcss", "eslint", "prettier"},
		},
		{
			name: "Vitest only, no Tailwind",
			cfg: models.Config{
				Styling: models.StylingCSSModules,
				Testing: models.TestingVitest,
			},
			wantOns: []string{"vitest", "eslint", "prettier"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildAddOns(tt.cfg)
			assertArgsEqual(t, got, tt.wantOns)
		})
	}
}

func TestBuildAddOns_AlwaysIncludesEslintPrettier(t *testing.T) {
	// Regardless of config, eslint and prettier must be present.
	cfg := models.Config{}
	got := buildAddOns(cfg)

	foundEslint := false
	foundPrettier := false
	for _, v := range got {
		if v == "eslint" {
			foundEslint = true
		}
		if v == "prettier" {
			foundPrettier = true
		}
	}
	if !foundEslint {
		t.Error("expected eslint in add-ons, not found")
	}
	if !foundPrettier {
		t.Error("expected prettier in add-ons, not found")
	}
}

func TestSupportedOptions(t *testing.T) {
	g := &Generator{}
	opts := g.SupportedOptions()

	t.Run("Styling", func(t *testing.T) {
		want := []string{"Tailwind CSS", "CSS Modules", "Sass/SCSS", "Vanilla CSS"}
		assertArgsEqual(t, opts.Styling, want)
	})

	t.Run("Testing", func(t *testing.T) {
		want := []string{"Vitest", "Playwright", "None"}
		assertArgsEqual(t, opts.Testing, want)
	})

	t.Run("StateManagement", func(t *testing.T) {
		want := []string{"Svelte Stores", "None"}
		assertArgsEqual(t, opts.StateManagement, want)
	})

	t.Run("DataFetching", func(t *testing.T) {
		want := []string{"TanStack Query", "Fetch API", "None"}
		assertArgsEqual(t, opts.DataFetching, want)
	})
}

func assertArgsEqual(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d\ngot:  %v\nwant: %v",
			len(got), len(want), got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got[i], want[i])
		}
	}
}
