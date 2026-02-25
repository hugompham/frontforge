package nextjs

import (
	"frontforge/internal/generators/meta"
	"frontforge/internal/models"
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
			name: "TypeScript with Tailwind and npm",
			cfg: models.Config{
				ProjectPath:    "/tmp/myapp",
				Language:       models.LangTypeScript,
				Styling:        models.StylingTailwind,
				PackageManager: models.PackageManagerNpm,
			},
			wantArgs: []string{
				"create-next-app@latest", "/tmp/myapp",
				"--ts", "--tailwind",
				"--eslint", "--app", "--src-dir", "--import-alias", "@/*", "--turbopack",
				"--use-npm", "--yes",
			},
		},
		{
			name: "JavaScript without Tailwind and yarn",
			cfg: models.Config{
				ProjectPath:    "/tmp/jsapp",
				Language:       models.LangJavaScript,
				Styling:        models.StylingVanilla,
				PackageManager: models.PackageManagerYarn,
			},
			wantArgs: []string{
				"create-next-app@latest", "/tmp/jsapp",
				"--js", "--no-tailwind",
				"--eslint", "--app", "--src-dir", "--import-alias", "@/*", "--turbopack",
				"--use-yarn", "--yes",
			},
		},
		{
			name: "TypeScript with CSS Modules and pnpm",
			cfg: models.Config{
				ProjectPath:    "/tmp/pnpmapp",
				Language:       models.LangTypeScript,
				Styling:        models.StylingCSSModules,
				PackageManager: models.PackageManagerPnpm,
			},
			wantArgs: []string{
				"create-next-app@latest", "/tmp/pnpmapp",
				"--ts", "--no-tailwind",
				"--eslint", "--app", "--src-dir", "--import-alias", "@/*", "--turbopack",
				"--use-pnpm", "--yes",
			},
		},
		{
			name: "JavaScript with Tailwind and bun",
			cfg: models.Config{
				ProjectPath:    "/tmp/bunapp",
				Language:       models.LangJavaScript,
				Styling:        models.StylingTailwind,
				PackageManager: models.PackageManagerBun,
			},
			wantArgs: []string{
				"create-next-app@latest", "/tmp/bunapp",
				"--js", "--tailwind",
				"--eslint", "--app", "--src-dir", "--import-alias", "@/*", "--turbopack",
				"--use-bun", "--yes",
			},
		},
		{
			name: "default package manager falls back to npm",
			cfg: models.Config{
				ProjectPath:    "/tmp/defaultpm",
				Language:       models.LangTypeScript,
				Styling:        models.StylingVanilla,
				PackageManager: "",
			},
			wantArgs: []string{
				"create-next-app@latest", "/tmp/defaultpm",
				"--ts", "--no-tailwind",
				"--eslint", "--app", "--src-dir", "--import-alias", "@/*", "--turbopack",
				"--use-npm", "--yes",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildScaffoldArgs(tt.cfg)
			if len(got) != len(tt.wantArgs) {
				t.Fatalf("arg count mismatch: got %d, want %d\ngot:  %v\nwant: %v",
					len(got), len(tt.wantArgs), got, tt.wantArgs)
			}
			for i := range got {
				if got[i] != tt.wantArgs[i] {
					t.Errorf("arg[%d] mismatch: got %q, want %q\nfull got:  %v\nfull want: %v",
						i, got[i], tt.wantArgs[i], got, tt.wantArgs)
				}
			}
		})
	}
}

func TestSupportedOptions(t *testing.T) {
	g := &Generator{}
	opts := g.SupportedOptions()

	t.Run("Styling options present", func(t *testing.T) {
		want := []string{"Tailwind CSS", "CSS Modules", "Sass/SCSS", "Vanilla CSS"}
		assertSliceEqual(t, opts.Styling, want)
	})

	t.Run("Testing options present", func(t *testing.T) {
		want := []string{"Vitest", "Jest", "None"}
		assertSliceEqual(t, opts.Testing, want)
	})

	t.Run("StateManagement options present", func(t *testing.T) {
		want := []string{"Zustand", "Redux Toolkit", "Context API", "None"}
		assertSliceEqual(t, opts.StateManagement, want)
	})

	t.Run("DataFetching options present", func(t *testing.T) {
		want := []string{"TanStack Query", "SWR", "Axios", "Fetch API", "None"}
		assertSliceEqual(t, opts.DataFetching, want)
	})
}

func assertSliceEqual(t *testing.T, got, want []string) {
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
