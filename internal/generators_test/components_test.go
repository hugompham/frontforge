package generators_test

import (
	"frontforge/internal/generators"
	"frontforge/internal/models"
	"strings"
	"testing"
)

func TestGenerateIndexHTML(t *testing.T) {
	tests := []struct {
		name          string
		config        models.Config
		wantMountID   string
		wantScriptExt string
	}{
		{
			name: "React TypeScript",
			config: models.Config{
				ProjectName: "react-app",
				Framework:   models.FrameworkReact,
				Language:    models.LangTypeScript,
			},
			wantMountID:   "root",
			wantScriptExt: "tsx",
		},
		{
			name: "Vue TypeScript",
			config: models.Config{
				ProjectName: "vue-app",
				Framework:   models.FrameworkVue,
				Language:    models.LangTypeScript,
			},
			wantMountID:   "app",
			wantScriptExt: "ts",
		},
		{
			name: "Svelte TypeScript",
			config: models.Config{
				ProjectName: "svelte-app",
				Framework:   models.FrameworkSvelte,
				Language:    models.LangTypeScript,
			},
			wantMountID:   "root",
			wantScriptExt: "ts",
		},
		{
			name: "Solid TypeScript",
			config: models.Config{
				ProjectName: "solid-app",
				Framework:   models.FrameworkSolid,
				Language:    models.LangTypeScript,
			},
			wantMountID:   "root",
			wantScriptExt: "tsx",
		},
		{
			name: "Angular TypeScript",
			config: models.Config{
				ProjectName: "angular-app",
				Framework:   models.FrameworkAngular,
				Language:    models.LangTypeScript,
			},
			wantMountID:   "root",
			wantScriptExt: "ts",
		},
		{
			name: "Vanilla TypeScript",
			config: models.Config{
				ProjectName: "vanilla-app",
				Framework:   models.FrameworkVanilla,
				Language:    models.LangTypeScript,
			},
			wantMountID:   "root",
			wantScriptExt: "ts",
		},
		{
			name: "Vanilla JavaScript",
			config: models.Config{
				ProjectName: "vanilla-js-app",
				Framework:   models.FrameworkVanilla,
				Language:    models.LangJavaScript,
			},
			wantMountID:   "root",
			wantScriptExt: "js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := generators.GenerateIndexHTML(tt.config)

			// Verify mount ID
			expectedMountID := `<div id="` + tt.wantMountID + `"></div>`
			if !strings.Contains(html, expectedMountID) {
				t.Errorf("Expected mount ID '%s' not found in HTML", tt.wantMountID)
			}

			// Verify script src extension
			expectedScript := `/src/main.` + tt.wantScriptExt
			if !strings.Contains(html, expectedScript) {
				t.Errorf("Expected script src '%s' not found. HTML:\n%s", expectedScript, html)
			}

			// Verify basic HTML structure
			if !strings.Contains(html, "<!doctype html>") {
				t.Error("Missing <!doctype html>")
			}
			if !strings.Contains(html, tt.config.ProjectName) {
				t.Errorf("Project name '%s' not found in HTML", tt.config.ProjectName)
			}
		})
	}
}

func TestGenerateMainFile(t *testing.T) {
	tests := []struct {
		name           string
		config         models.Config
		wantImports    []string
		wantNotContain string
	}{
		{
			name: "React",
			config: models.Config{
				Framework: models.FrameworkReact,
				Language:  models.LangTypeScript,
			},
			wantImports: []string{"from 'react'", "from 'react-dom/client'", "createRoot"},
		},
		{
			name: "Vue",
			config: models.Config{
				Framework: models.FrameworkVue,
				Language:  models.LangTypeScript,
			},
			wantImports: []string{"from 'vue'", "createApp", "App.vue"},
		},
		{
			name: "Angular",
			config: models.Config{
				Framework: models.FrameworkAngular,
				Language:  models.LangTypeScript,
			},
			wantImports: []string{"bootstrapApplication", "@angular/platform-browser", "AppComponent"},
		},
		{
			name: "Svelte",
			config: models.Config{
				Framework: models.FrameworkSvelte,
				Language:  models.LangTypeScript,
			},
			wantImports: []string{"from 'svelte'", "mount", "App.svelte"},
		},
		{
			name: "Solid",
			config: models.Config{
				Framework: models.FrameworkSolid,
				Language:  models.LangTypeScript,
			},
			wantImports: []string{"from 'solid-js/web'", "render"},
		},
		{
			name: "Vanilla",
			config: models.Config{
				Framework: models.FrameworkVanilla,
				Language:  models.LangTypeScript,
			},
			wantImports: []string{"import App", "getElementById"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mainFile := generators.GenerateMainFile(tt.config)

			for _, imp := range tt.wantImports {
				if !strings.Contains(mainFile, imp) {
					t.Errorf("Expected import/code '%s' not found in main file", imp)
				}
			}

			if tt.wantNotContain != "" && strings.Contains(mainFile, tt.wantNotContain) {
				t.Errorf("Should not contain '%s' in main file", tt.wantNotContain)
			}
		})
	}
}

func TestGenerateAppFile(t *testing.T) {
	tests := []struct {
		name        string
		config      models.Config
		wantContent []string
	}{
		{
			name: "React TypeScript",
			config: models.Config{
				ProjectName: "react-app",
				Framework:   models.FrameworkReact,
				Language:    models.LangTypeScript,
			},
			wantContent: []string{"useState", "function App()", "export default App"},
		},
		{
			name: "Vue",
			config: models.Config{
				ProjectName: "vue-app",
				Framework:   models.FrameworkVue,
				Language:    models.LangTypeScript,
			},
			wantContent: []string{"<script setup", "from 'vue'", "<template>", "<style scoped>"},
		},
		{
			name: "Angular",
			config: models.Config{
				ProjectName: "angular-app",
				Framework:   models.FrameworkAngular,
				Language:    models.LangTypeScript,
			},
			wantContent: []string{"@Component", "standalone: true", "export class AppComponent"},
		},
		{
			name: "Svelte",
			config: models.Config{
				ProjectName: "svelte-app",
				Framework:   models.FrameworkSvelte,
				Language:    models.LangTypeScript,
			},
			wantContent: []string{"<script>", "$state", "<main>", "<style>"},
		},
		{
			name: "Solid",
			config: models.Config{
				ProjectName: "solid-app",
				Framework:   models.FrameworkSolid,
				Language:    models.LangTypeScript,
			},
			wantContent: []string{"createSignal", "function App()", "export default App"},
		},
		{
			name: "Vanilla TypeScript",
			config: models.Config{
				ProjectName: "vanilla-app",
				Framework:   models.FrameworkVanilla,
				Language:    models.LangTypeScript,
			},
			wantContent: []string{"export default function App()", "addEventListener"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appFile := generators.GenerateAppFile(tt.config)

			for _, want := range tt.wantContent {
				if !strings.Contains(appFile, want) {
					t.Errorf("Expected content '%s' not found in app file", want)
				}
			}

			// All apps should have project name
			if !strings.Contains(appFile, tt.config.ProjectName) {
				t.Errorf("Project name '%s' not found in app file", tt.config.ProjectName)
			}
		})
	}
}
