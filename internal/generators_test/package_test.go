package generators_test

import (
	"frontforge/internal/generators"
	"frontforge/internal/models"
	"testing"
)

func TestGeneratePackageJSON(t *testing.T) {
	tests := []struct {
		name     string
		config   models.Config
		wantName string
	}{
		{
			name: "React with TypeScript",
			config: models.Config{
				ProjectName: "test-project",
				Language:    models.LangTypeScript,
				Framework:   models.FrameworkReact,
			},
			wantName: "test-project",
		},
		{
			name: "Vue with JavaScript",
			config: models.Config{
				ProjectName: "vue-app",
				Language:    models.LangJavaScript,
				Framework:   models.FrameworkVue,
			},
			wantName: "vue-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := generators.GeneratePackageJSON(tt.config)

			if pkg.Name != tt.wantName {
				t.Errorf("GeneratePackageJSON() name = %v, want %v", pkg.Name, tt.wantName)
			}

			if !pkg.Private {
				t.Error("GeneratePackageJSON() should set Private = true")
			}

			if pkg.Type != "module" {
				t.Errorf("GeneratePackageJSON() type = %v, want module", pkg.Type)
			}

			// Check framework dependencies
			if tt.config.Framework == models.FrameworkReact {
				if _, exists := pkg.Dependencies["react"]; !exists {
					t.Error("React project should have react dependency")
				}
			}

			// Check TypeScript dependencies
			if tt.config.Language == models.LangTypeScript {
				if _, exists := pkg.DevDependencies["typescript"]; !exists {
					t.Error("TypeScript project should have typescript dev dependency")
				}
			}
		})
	}
}

func TestGeneratePackageJSON_Scripts(t *testing.T) {
	config := models.Config{
		ProjectName: "test-project",
		Framework:   models.FrameworkReact,
		Testing:     models.TestingVitest,
	}

	pkg := generators.GeneratePackageJSON(config)

	expectedScripts := []string{"dev", "build", "preview", "lint", "test"}
	for _, script := range expectedScripts {
		if _, exists := pkg.Scripts[script]; !exists {
			t.Errorf("Expected script '%s' not found", script)
		}
	}
}
