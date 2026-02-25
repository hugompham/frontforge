package astro

import (
	"frontforge/internal/generators/meta"
	"frontforge/internal/models"
	"strings"
	"testing"
)

// Compile-time interface compliance check.
var _ meta.MetaGenerator = (*Generator)(nil)

func TestBuildScaffoldArgs(t *testing.T) {
	tests := []struct {
		name     string
		cfg      models.Config
		wantArgs []string
	}{
		{
			name: "TypeScript strict",
			cfg: models.Config{
				ProjectPath: "/tmp/astro-ts",
				Language:    models.LangTypeScript,
			},
			wantArgs: []string{
				"create", "astro@latest", "/tmp/astro-ts", "--",
				"--template", "minimal",
				"--typescript", "strict",
				"--install", "--git", "--skip-houston",
			},
		},
		{
			name: "JavaScript relaxed",
			cfg: models.Config{
				ProjectPath: "/tmp/astro-js",
				Language:    models.LangJavaScript,
			},
			wantArgs: []string{
				"create", "astro@latest", "/tmp/astro-js", "--",
				"--template", "minimal",
				"--typescript", "relaxed",
				"--install", "--git", "--skip-houston",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildScaffoldArgs(tt.cfg)
			assertArgsEqual(t, got, tt.wantArgs)
		})
	}
}

func TestBuildScaffoldArgs_AlwaysIncludesRequiredFlags(t *testing.T) {
	// --skip-houston, --install, --git must always be present.
	cfg := models.Config{
		ProjectPath: "/tmp/test",
		Language:    models.LangTypeScript,
	}
	got := buildScaffoldArgs(cfg)
	joined := strings.Join(got, " ")

	for _, flag := range []string{"--skip-houston", "--install", "--git"} {
		if !strings.Contains(joined, flag) {
			t.Errorf("expected flag %q in args, not found: %v", flag, got)
		}
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
		want := []string{"Vitest", "None"}
		assertArgsEqual(t, opts.Testing, want)
	})

	t.Run("StateManagement is nil (hidden)", func(t *testing.T) {
		if opts.StateManagement != nil {
			t.Errorf("expected StateManagement to be nil, got %v", opts.StateManagement)
		}
	})

	t.Run("DataFetching is nil (hidden)", func(t *testing.T) {
		if opts.DataFetching != nil {
			t.Errorf("expected DataFetching to be nil, got %v", opts.DataFetching)
		}
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
