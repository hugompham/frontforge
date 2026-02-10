package integration_test

import (
	"frontforge/internal/generators"
	"frontforge/internal/models"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestEndToEndGeneration tests complete project generation
func TestEndToEndGeneration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name           string
		config         models.Config
		expectedFiles  []string
		expectedInFile map[string][]string // file -> content snippets
	}{
		{
			name: "React TypeScript with Tailwind",
			config: models.Config{
				ProjectName:     "test-react-app",
				Language:        models.LangTypeScript,
				Framework:       models.FrameworkReact,
				PackageManager:  models.PackageManagerNpm,
				Styling:         models.StylingTailwind,
				Routing:         models.RoutingReactRouter,
				Testing:         models.TestingVitest,
				StateManagement: models.StateZustand,
				DataFetching:    models.DataTanStackQuery,
				Structure:       models.StructureFeatureBased,
			},
			expectedFiles: []string{
				"package.json",
				"tsconfig.json",
				"vite.config.ts",
				"index.html",
				"eslint.config.js",
				"vitest.config.ts",
				".gitignore",
				"README.md",
				"src/main.tsx",
				"src/App.tsx",
				"src/index.css",
				"src/test/setup.ts",
				"public/vite.svg",
			},
			expectedInFile: map[string][]string{
				"package.json": {
					"\"name\": \"test-react-app\"",
					"\"react\":",
					"\"vite\":",
					"\"typescript\":",
				},
				"src/main.tsx": {
					"from 'react'",
					"import { createRoot }",
					"import App from",
				},
				"src/App.tsx": {
					"function App()",
					"export default App",
				},
				"README.md": {
					"# test-react-app",
					"React",
					"TypeScript",
					"Tailwind CSS",
				},
			},
		},
		{
			name: "Vue JavaScript with CSS Modules",
			config: models.Config{
				ProjectName:    "test-vue-app",
				Language:       models.LangJavaScript,
				Framework:      models.FrameworkVue,
				PackageManager: models.PackageManagerPnpm,
				Styling:        models.StylingCSSModules,
				Routing:        models.RoutingNone,
				Testing:        models.TestingNone,
				Structure:      models.StructureLayerBased,
			},
			expectedFiles: []string{
				"package.json",
				"vite.config.js",
				"index.html",
				"eslint.config.js",
				".gitignore",
				"README.md",
				"src/main.js",
				"src/App.vue",
				"src/App.module.css",
				"public/vite.svg",
			},
			expectedInFile: map[string][]string{
				"package.json": {
					"\"name\": \"test-vue-app\"",
					"\"vue\":",
				},
				"src/App.vue": {
					"<template>",
					"</template>",
					"<script setup",
					"</script>",
				},
				"README.md": {
					"# test-vue-app",
					"Vue",
					"JavaScript",
				},
			},
		},
		{
			name: "Vanilla TypeScript minimal",
			config: models.Config{
				ProjectName:    "vanilla-app",
				Language:       models.LangTypeScript,
				Framework:      models.FrameworkVanilla,
				PackageManager: models.PackageManagerBun,
				Styling:        "Vanilla CSS",
				Routing:        models.RoutingNone,
				Testing:        models.TestingNone,
				Structure:      models.StructureLayerBased,
			},
			expectedFiles: []string{
				"package.json",
				"tsconfig.json",
				"vite.config.ts",
				"index.html",
				"eslint.config.js",
				".gitignore",
				"README.md",
				"src/main.ts",
				"public/vite.svg",
			},
			expectedInFile: map[string][]string{
				"package.json": {
					"\"name\": \"vanilla-app\"",
					"\"vite\":",
				},
				"README.md": {
					"# vanilla-app",
					"Vanilla",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()
			projectPath := filepath.Join(tmpDir, tt.config.ProjectName)
			tt.config.ProjectPath = projectPath

			// Generate project
			err := generators.SetupProject(tt.config)
			if err != nil {
				t.Fatalf("SetupProject failed: %v", err)
			}

			// Verify all expected files exist
			for _, expectedFile := range tt.expectedFiles {
				fullPath := filepath.Join(projectPath, expectedFile)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("Expected file not found: %s", expectedFile)
				}
			}

			// Verify file contents
			for file, expectedSnippets := range tt.expectedInFile {
				fullPath := filepath.Join(projectPath, file)
				content, err := os.ReadFile(fullPath)
				if err != nil {
					t.Errorf("Failed to read %s: %v", file, err)
					continue
				}

				contentStr := string(content)
				for _, snippet := range expectedSnippets {
					if !contains(contentStr, snippet) {
						t.Errorf("File %s missing expected content: %q", file, snippet)
					}
				}
			}
		})
	}
}

// TestDryRunMode tests that dry-run mode doesn't create files
func TestDryRunMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "dry-run-test")

	config := models.Config{
		ProjectName:    "dry-run-test",
		ProjectPath:    projectPath,
		Language:       models.LangTypeScript,
		Framework:      models.FrameworkReact,
		PackageManager: models.PackageManagerNpm,
		Styling:        models.StylingTailwind,
		DryRun:         true,
	}

	err := generators.SetupProject(config)
	if err != nil {
		t.Fatalf("SetupProject failed in dry-run mode: %v", err)
	}

	// Verify project directory was NOT created
	if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
		t.Error("Dry-run mode should not create project directory")
	}
}

// TestValidateProject tests post-generation validation
func TestValidateProject(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "validate-test")

	config := models.Config{
		ProjectName:    "validate-test",
		ProjectPath:    projectPath,
		Language:       models.LangTypeScript,
		Framework:      models.FrameworkReact,
		PackageManager: models.PackageManagerNpm,
		Styling:        models.StylingTailwind,
		Testing:        models.TestingVitest,
	}

	// Generate project
	err := generators.SetupProject(config)
	if err != nil {
		t.Fatalf("SetupProject failed: %v", err)
	}

	// Run validation
	results := generators.ValidateProject(projectPath, config)

	// Check that validation ran
	if len(results) == 0 {
		t.Error("Expected validation results, got none")
	}

	// All critical files should pass validation
	criticalChecks := []string{
		"File exists: package.json",
		"File exists: index.html",
		"vite.config exists",
	}

	for _, check := range criticalChecks {
		found := false
		for _, result := range results {
			if result.Check == check && result.Passed {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Critical check failed or missing: %s", check)
		}
	}
}

// TestInvalidProjectPath tests error handling for invalid paths
func TestInvalidProjectPath(t *testing.T) {
	// Use platform-specific forbidden paths
	var tests []struct {
		name        string
		projectPath string
		wantError   bool
	}

	if runtime.GOOS == "windows" {
		tests = []struct {
			name        string
			projectPath string
			wantError   bool
		}{
			{
				name:        "Windows system directory",
				projectPath: "C:\\Windows",
				wantError:   true,
			},
			{
				name:        "Program Files",
				projectPath: "C:\\Program Files",
				wantError:   true,
			},
		}
	} else {
		// Unix-like systems (Linux, macOS)
		tests = []struct {
			name        string
			projectPath string
			wantError   bool
		}{
			{
				name:        "Root directory",
				projectPath: "/",
				wantError:   true,
			},
			{
				name:        "System directory",
				projectPath: "/usr/bin/test",
				wantError:   true,
			},
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := models.Config{
				ProjectName:    "test",
				ProjectPath:    tt.projectPath,
				Language:       models.LangTypeScript,
				Framework:      models.FrameworkReact,
				PackageManager: models.PackageManagerNpm,
			}

			err := generators.SetupProject(config)
			if tt.wantError && err == nil {
				t.Error("Expected error for invalid path, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
