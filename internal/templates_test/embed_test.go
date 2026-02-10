package templates_test

import (
	"frontforge/internal/models"
	"frontforge/internal/templates"
	"strings"
	"testing"
)

func TestPrepareTemplateData_ComputedFields(t *testing.T) {
	tests := []struct {
		name                string
		config              models.Config
		wantMountID         string
		wantMainExt         string
		wantAppExt          string
		wantPmRun           string
		wantVitestExt       string
		wantFrameworkDocURL string
	}{
		{
			name: "React TypeScript with npm",
			config: models.Config{
				ProjectName:    "test-app",
				Framework:      models.FrameworkReact,
				Language:       models.LangTypeScript,
				PackageManager: models.PackageManagerNpm,
			},
			wantMountID:         "root",
			wantMainExt:         "tsx",
			wantAppExt:          "tsx",
			wantPmRun:           "npm run",
			wantVitestExt:       "ts",
			wantFrameworkDocURL: "https://react.dev",
		},
		{
			name: "Vue TypeScript with pnpm",
			config: models.Config{
				ProjectName:    "vue-app",
				Framework:      models.FrameworkVue,
				Language:       models.LangTypeScript,
				PackageManager: models.PackageManagerPnpm,
			},
			wantMountID:         "app",
			wantMainExt:         "ts",
			wantAppExt:          "vue",
			wantPmRun:           "pnpm",
			wantVitestExt:       "ts",
			wantFrameworkDocURL: "https://vuejs.org",
		},
		{
			name: "Svelte JavaScript with yarn",
			config: models.Config{
				ProjectName:    "svelte-app",
				Framework:      models.FrameworkSvelte,
				Language:       models.LangJavaScript,
				PackageManager: models.PackageManagerYarn,
			},
			wantMountID:         "root",
			wantMainExt:         "js",
			wantAppExt:          "svelte",
			wantPmRun:           "yarn",
			wantVitestExt:       "js",
			wantFrameworkDocURL: "https://svelte.dev",
		},
		{
			name: "Solid TypeScript with bun",
			config: models.Config{
				ProjectName:    "solid-app",
				Framework:      models.FrameworkSolid,
				Language:       models.LangTypeScript,
				PackageManager: models.PackageManagerBun,
			},
			wantMountID:         "root",
			wantMainExt:         "tsx",
			wantAppExt:          "tsx",
			wantPmRun:           "bun",
			wantVitestExt:       "ts",
			wantFrameworkDocURL: "https://www.solidjs.com",
		},
		{
			name: "Angular TypeScript",
			config: models.Config{
				ProjectName:    "ng-app",
				Framework:      models.FrameworkAngular,
				Language:       models.LangTypeScript,
				PackageManager: models.PackageManagerNpm,
			},
			wantMountID:         "root",
			wantMainExt:         "ts",
			wantAppExt:          "ts",
			wantPmRun:           "npm run",
			wantVitestExt:       "ts",
			wantFrameworkDocURL: "https://angular.dev",
		},
		{
			name: "Vanilla JavaScript",
			config: models.Config{
				ProjectName:    "vanilla-app",
				Framework:      models.FrameworkVanilla,
				Language:       models.LangJavaScript,
				PackageManager: models.PackageManagerNpm,
			},
			wantMountID:         "root",
			wantMainExt:         "js",
			wantAppExt:          "jsx",
			wantPmRun:           "npm run",
			wantVitestExt:       "js",
			wantFrameworkDocURL: "https://developer.mozilla.org/en-US/docs/Web/JavaScript",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := templates.PrepareTemplateData(tt.config)

			if data.MountID != tt.wantMountID {
				t.Errorf("MountID = %q, want %q", data.MountID, tt.wantMountID)
			}
			if data.MainExt != tt.wantMainExt {
				t.Errorf("MainExt = %q, want %q", data.MainExt, tt.wantMainExt)
			}
			if data.AppExt != tt.wantAppExt {
				t.Errorf("AppExt = %q, want %q", data.AppExt, tt.wantAppExt)
			}
			if data.PmRun != tt.wantPmRun {
				t.Errorf("PmRun = %q, want %q", data.PmRun, tt.wantPmRun)
			}
			if data.VitestExt != tt.wantVitestExt {
				t.Errorf("VitestExt = %q, want %q", data.VitestExt, tt.wantVitestExt)
			}
			if data.FrameworkDocURL != tt.wantFrameworkDocURL {
				t.Errorf("FrameworkDocURL = %q, want %q", data.FrameworkDocURL, tt.wantFrameworkDocURL)
			}
		})
	}
}

func TestPrepareTemplateData_StructureExample(t *testing.T) {
	tests := []struct {
		name            string
		config          models.Config
		wantContains    []string
		wantNotContains []string
	}{
		{
			name: "Feature-based TypeScript",
			config: models.Config{
				Structure: models.StructureFeatureBased,
				Language:  models.LangTypeScript,
			},
			wantContains:    []string{"features/", "auth/", "dashboard/", "App.tsx", "main.tsx"},
			wantNotContains: []string{"pages/", "services/"},
		},
		{
			name: "Layer-based JavaScript",
			config: models.Config{
				Structure: models.StructureLayerBased,
				Language:  models.LangJavaScript,
			},
			wantContains:    []string{"components/", "pages/", "services/", "App.jsx", "main.jsx"},
			wantNotContains: []string{"features/", "auth/", "dashboard/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := templates.PrepareTemplateData(tt.config)

			for _, want := range tt.wantContains {
				if !strings.Contains(data.StructureExample, want) {
					t.Errorf("StructureExample missing %q", want)
				}
			}
			for _, notWant := range tt.wantNotContains {
				if strings.Contains(data.StructureExample, notWant) {
					t.Errorf("StructureExample should not contain %q", notWant)
				}
			}
		})
	}
}

func TestRenderVitestConfig(t *testing.T) {
	tests := []struct {
		name         string
		framework    string
		wantContains []string
	}{
		{
			name:         "React",
			framework:    models.FrameworkReact,
			wantContains: []string{"import react from", "plugins: [react()"},
		},
		{
			name:         "Vue",
			framework:    models.FrameworkVue,
			wantContains: []string{"import vue from", "plugins: [vue()"},
		},
		{
			name:         "Svelte",
			framework:    models.FrameworkSvelte,
			wantContains: []string{"import { svelte }", "plugins: [svelte()"},
		},
		{
			name:         "Vanilla",
			framework:    models.FrameworkVanilla,
			wantContains: []string{"plugins: []"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := models.Config{
				Framework: tt.framework,
				Language:  models.LangTypeScript,
			}

			result, err := templates.RenderVitestConfig(config)
			if err != nil {
				t.Fatalf("RenderVitestConfig() error = %v", err)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("result missing %q", want)
				}
			}
		})
	}
}

func TestRenderVitestSetup(t *testing.T) {
	tests := []struct {
		name         string
		framework    string
		wantContains []string
	}{
		{
			name:         "React",
			framework:    models.FrameworkReact,
			wantContains: []string{"import { cleanup } from '@testing-library/react'", "afterEach"},
		},
		{
			name:         "Vue",
			framework:    models.FrameworkVue,
			wantContains: []string{"import { expect } from 'vitest'", "@testing-library/jest-dom"},
		},
		{
			name:         "Svelte",
			framework:    models.FrameworkSvelte,
			wantContains: []string{"import { cleanup } from '@testing-library/svelte'", "afterEach"},
		},
		{
			name:         "Vanilla",
			framework:    models.FrameworkVanilla,
			wantContains: []string{"import { expect } from 'vitest'", "@testing-library/jest-dom"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := models.Config{
				Framework: tt.framework,
				Language:  models.LangTypeScript,
			}

			result, err := templates.RenderVitestSetup(config)
			if err != nil {
				t.Fatalf("RenderVitestSetup() error = %v", err)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("result missing %q", want)
				}
			}
		})
	}
}

func TestRenderESLintConfig(t *testing.T) {
	tests := []struct {
		name         string
		framework    string
		wantContains []string
	}{
		{
			name:         "React",
			framework:    models.FrameworkReact,
			wantContains: []string{"eslint-plugin-react-hooks", "react-refresh"},
		},
		{
			name:         "Vue",
			framework:    models.FrameworkVue,
			wantContains: []string{"eslint-plugin-vue"},
		},
		{
			name:         "Svelte",
			framework:    models.FrameworkSvelte,
			wantContains: []string{"typescript-eslint"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := models.Config{
				Framework: tt.framework,
			}

			result, err := templates.RenderESLintConfig(config)
			if err != nil {
				t.Fatalf("RenderESLintConfig() error = %v", err)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("result missing %q", want)
				}
			}
		})
	}
}

func TestRenderREADME(t *testing.T) {
	config := models.Config{
		ProjectName:     "test-project",
		Framework:       models.FrameworkReact,
		Language:        models.LangTypeScript,
		PackageManager:  models.PackageManagerNpm,
		Styling:         models.StylingTailwind,
		Routing:         models.RoutingReactRouter,
		Testing:         models.TestingVitest,
		StateManagement: "Zustand",
		DataFetching:    models.DataTanStackQuery,
		Structure:       models.StructureFeatureBased,
	}

	result, err := templates.Render("static/README.md.tmpl", config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	wantContains := []string{
		"# test-project",
		"**Language**: TypeScript",
		"**Framework**: React",
		"npm run dev",
		"https://react.dev",
	}

	for _, want := range wantContains {
		if !strings.Contains(result, want) {
			t.Errorf("README missing %q", want)
		}
	}
}
