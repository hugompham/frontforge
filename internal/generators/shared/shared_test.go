package shared

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// --- helpers ---

func writePackageJSON(t *testing.T, dir string, pkg map[string]interface{}) {
	t.Helper()
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		t.Fatalf("marshal package.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), data, 0644); err != nil {
		t.Fatalf("write package.json: %v", err)
	}
}

func readPackageJSON(t *testing.T, dir string) map[string]interface{} {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatalf("read package.json: %v", err)
	}
	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		t.Fatalf("parse package.json: %v", err)
	}
	return pkg
}

func getStringMap(pkg map[string]interface{}, key string) map[string]string {
	raw, ok := pkg[key].(map[string]interface{})
	if !ok {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		out[k] = v.(string)
	}
	return out
}

func dirExists(t *testing.T, path string) bool {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func fileExists(t *testing.T, path string) bool {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// --- MergePackageJSON ---

func TestMergePackageJSON(t *testing.T) {
	tests := []struct {
		name        string
		existing    map[string]interface{}
		deps        map[string]string
		devDeps     map[string]string
		scripts     map[string]string
		wantDeps    map[string]string
		wantDevDeps map[string]string
		wantScripts map[string]string
	}{
		{
			name: "adds new deps to empty package.json",
			existing: map[string]interface{}{
				"name": "test-project",
			},
			deps:        map[string]string{"react": "^19.0.0"},
			devDeps:     map[string]string{"vitest": "^4.0.0"},
			scripts:     map[string]string{"test": "vitest"},
			wantDeps:    map[string]string{"react": "^19.0.0"},
			wantDevDeps: map[string]string{"vitest": "^4.0.0"},
			wantScripts: map[string]string{"test": "vitest"},
		},
		{
			name: "preserves existing deps, does not overwrite",
			existing: map[string]interface{}{
				"name": "test-project",
				"dependencies": map[string]interface{}{
					"react": "^18.0.0",
				},
				"devDependencies": map[string]interface{}{
					"vitest": "^3.0.0",
				},
				"scripts": map[string]interface{}{
					"dev": "next dev",
				},
			},
			deps:        map[string]string{"react": "^19.0.0", "zustand": "^5.0.0"},
			devDeps:     map[string]string{"vitest": "^4.0.0", "jsdom": "^25.0.0"},
			scripts:     map[string]string{"dev": "vite", "test": "vitest"},
			wantDeps:    map[string]string{"react": "^18.0.0", "zustand": "^5.0.0"},
			wantDevDeps: map[string]string{"vitest": "^3.0.0", "jsdom": "^25.0.0"},
			wantScripts: map[string]string{"dev": "next dev", "test": "vitest"},
		},
		{
			name: "nil maps are no-ops",
			existing: map[string]interface{}{
				"name": "test-project",
				"dependencies": map[string]interface{}{
					"react": "^18.0.0",
				},
			},
			deps:     nil,
			devDeps:  nil,
			scripts:  nil,
			wantDeps: map[string]string{"react": "^18.0.0"},
		},
		{
			name: "empty maps are no-ops",
			existing: map[string]interface{}{
				"name": "test-project",
			},
			deps:    map[string]string{},
			devDeps: map[string]string{},
			scripts: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writePackageJSON(t, dir, tt.existing)

			err := MergePackageJSON(dir, tt.deps, tt.devDeps, tt.scripts)
			if err != nil {
				t.Fatalf("MergePackageJSON returned error: %v", err)
			}

			pkg := readPackageJSON(t, dir)

			if tt.wantDeps != nil {
				got := getStringMap(pkg, "dependencies")
				for k, want := range tt.wantDeps {
					if got[k] != want {
						t.Errorf("dependencies[%q] = %q, want %q", k, got[k], want)
					}
				}
			}

			if tt.wantDevDeps != nil {
				got := getStringMap(pkg, "devDependencies")
				for k, want := range tt.wantDevDeps {
					if got[k] != want {
						t.Errorf("devDependencies[%q] = %q, want %q", k, got[k], want)
					}
				}
			}

			if tt.wantScripts != nil {
				got := getStringMap(pkg, "scripts")
				for k, want := range tt.wantScripts {
					if got[k] != want {
						t.Errorf("scripts[%q] = %q, want %q", k, got[k], want)
					}
				}
			}
		})
	}
}

func TestMergePackageJSON_MissingFile(t *testing.T) {
	dir := t.TempDir()
	err := MergePackageJSON(dir, map[string]string{"react": "^19.0.0"}, nil, nil)
	if err == nil {
		t.Fatal("expected error for missing package.json, got nil")
	}
}

func TestMergePackageJSON_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("{invalid"), 0644); err != nil {
		t.Fatal(err)
	}
	err := MergePackageJSON(dir, map[string]string{"react": "^19.0.0"}, nil, nil)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// --- AddNpmScripts ---

func TestAddNpmScripts(t *testing.T) {
	tests := []struct {
		name        string
		existing    map[string]interface{}
		scripts     map[string]string
		wantScripts map[string]string
	}{
		{
			name: "adds scripts to package.json without scripts field",
			existing: map[string]interface{}{
				"name": "test-project",
			},
			scripts:     map[string]string{"test": "vitest", "build": "next build"},
			wantScripts: map[string]string{"test": "vitest", "build": "next build"},
		},
		{
			name: "preserves existing scripts",
			existing: map[string]interface{}{
				"name": "test-project",
				"scripts": map[string]interface{}{
					"dev": "next dev",
				},
			},
			scripts:     map[string]string{"dev": "vite dev", "test": "vitest"},
			wantScripts: map[string]string{"dev": "next dev", "test": "vitest"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writePackageJSON(t, dir, tt.existing)

			err := AddNpmScripts(dir, tt.scripts)
			if err != nil {
				t.Fatalf("AddNpmScripts returned error: %v", err)
			}

			pkg := readPackageJSON(t, dir)
			got := getStringMap(pkg, "scripts")

			for k, want := range tt.wantScripts {
				if got[k] != want {
					t.Errorf("scripts[%q] = %q, want %q", k, got[k], want)
				}
			}
		})
	}
}

// --- ScaffoldVitest ---

func TestScaffoldVitest(t *testing.T) {
	tests := []struct {
		name          string
		framework     string
		wantDevDeps   []string
		wantNoDevDeps []string
		wantConfig    string
	}{
		{
			name:      "nextjs",
			framework: "nextjs",
			wantDevDeps: []string{
				"vitest", "@testing-library/jest-dom", "jsdom",
				"@testing-library/react", "@vitejs/plugin-react",
			},
			wantConfig: "vitest.config.ts",
		},
		{
			name:      "sveltekit",
			framework: "sveltekit",
			wantDevDeps: []string{
				"vitest", "@testing-library/jest-dom", "jsdom",
				"@testing-library/svelte",
			},
			wantNoDevDeps: []string{"@testing-library/react"},
			wantConfig:    "vitest.config.ts",
		},
		{
			name:      "astro",
			framework: "astro",
			wantDevDeps: []string{
				"vitest", "@testing-library/jest-dom", "jsdom",
			},
			wantNoDevDeps: []string{"@testing-library/react", "@testing-library/svelte", "@vitejs/plugin-react"},
			wantConfig:    "vitest.config.ts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writePackageJSON(t, dir, map[string]interface{}{
				"name": "test-" + tt.framework,
			})

			err := ScaffoldVitest(dir, tt.framework)
			if err != nil {
				t.Fatalf("ScaffoldVitest returned error: %v", err)
			}

			// Verify config file created
			configPath := filepath.Join(dir, tt.wantConfig)
			if !fileExists(t, configPath) {
				t.Errorf("expected config file %s to exist", tt.wantConfig)
			}

			// Verify test setup file created
			setupPath := filepath.Join(dir, "src", "test", "setup.ts")
			if !fileExists(t, setupPath) {
				t.Errorf("expected test setup file at src/test/setup.ts")
			}

			// Verify devDependencies merged
			pkg := readPackageJSON(t, dir)
			devDeps := getStringMap(pkg, "devDependencies")

			for _, dep := range tt.wantDevDeps {
				if devDeps[dep] == "" {
					t.Errorf("expected devDependency %q to be present", dep)
				}
			}

			for _, dep := range tt.wantNoDevDeps {
				if devDeps[dep] != "" {
					t.Errorf("unexpected devDependency %q present", dep)
				}
			}

			// Verify test script added
			scripts := getStringMap(pkg, "scripts")
			if scripts["test"] != "vitest" {
				t.Errorf("scripts[test] = %q, want %q", scripts["test"], "vitest")
			}
		})
	}
}

// --- ScaffoldFeatureStructure ---

func TestScaffoldFeatureStructure(t *testing.T) {
	tests := []struct {
		name      string
		framework string
		wantDirs  []string
	}{
		{
			name:      "nextjs creates app-based dirs",
			framework: "nextjs",
			wantDirs: []string{
				filepath.Join("app", "features"),
				filepath.Join("app", "components"),
				"lib",
				"hooks",
			},
		},
		{
			name:      "sveltekit creates src/lib dirs",
			framework: "sveltekit",
			wantDirs: []string{
				filepath.Join("src", "lib", "features"),
				filepath.Join("src", "lib", "components"),
				filepath.Join("src", "lib", "stores"),
			},
		},
		{
			name:      "astro creates src dirs",
			framework: "astro",
			wantDirs: []string{
				filepath.Join("src", "components"),
				filepath.Join("src", "layouts"),
				filepath.Join("src", "styles"),
			},
		},
		{
			name:      "unknown framework creates nothing",
			framework: "unknown",
			wantDirs:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			err := ScaffoldFeatureStructure(dir, tt.framework)
			if err != nil {
				t.Fatalf("ScaffoldFeatureStructure returned error: %v", err)
			}

			for _, rel := range tt.wantDirs {
				full := filepath.Join(dir, rel)
				if !dirExists(t, full) {
					t.Errorf("expected directory %s to exist", rel)
				}
			}
		})
	}
}
