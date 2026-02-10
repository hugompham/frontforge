package generators_test

import (
	"frontforge/internal/generators"
	"frontforge/internal/models"
	"frontforge/internal/testutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateViteConfig(t *testing.T) {
	tests := []struct {
		name         string
		config       models.Config
		wantPlugins  []string
	}{
		{
			name: "React with Tailwind",
			config: models.Config{
				Framework: models.FrameworkReact,
				Styling:   models.StylingTailwind,
			},
			wantPlugins: []string{"@vitejs/plugin-react", "react()", "@tailwindcss/vite", "tailwindcss()"},
		},
		{
			name: "Vue",
			config: models.Config{
				Framework: models.FrameworkVue,
				Styling:   models.StylingVanilla,
			},
			wantPlugins: []string{"@vitejs/plugin-vue", "vue()"},
		},
		{
			name: "Svelte",
			config: models.Config{
				Framework: models.FrameworkSvelte,
				Styling:   models.StylingVanilla,
			},
			wantPlugins: []string{"@sveltejs/vite-plugin-svelte", "svelte()"},
		},
		{
			name: "Solid",
			config: models.Config{
				Framework: models.FrameworkSolid,
				Styling:   models.StylingVanilla,
			},
			wantPlugins: []string{"vite-plugin-solid", "solid()"},
		},
		{
			name: "Angular",
			config: models.Config{
				Framework: models.FrameworkAngular,
				Styling:   models.StylingVanilla,
			},
			wantPlugins: []string{"@analogjs/vite-plugin-angular", "angular()"},
		},
		{
			name: "Vanilla",
			config: models.Config{
				Framework: models.FrameworkVanilla,
				Styling:   models.StylingVanilla,
			},
			wantPlugins: []string{"defineConfig", "plugins: []"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viteConfig := generators.GenerateViteConfig(tt.config)

			for _, plugin := range tt.wantPlugins {
				if !strings.Contains(viteConfig, plugin) {
					t.Errorf("Expected plugin '%s' not found in vite config", plugin)
				}
			}

			// All configs should have defineConfig
			if !strings.Contains(viteConfig, "defineConfig") {
				t.Error("Missing defineConfig import")
			}
		})
	}
}

func TestGenerateTSConfig(t *testing.T) {
	tests := []struct {
		name              string
		config            models.Config
		wantJsx           string
		wantExtraOptions  map[string]interface{}
	}{
		{
			name: "React",
			config: models.Config{
				Framework: models.FrameworkReact,
			},
			wantJsx: "react-jsx",
		},
		{
			name: "Vue",
			config: models.Config{
				Framework: models.FrameworkVue,
			},
			wantJsx: "preserve",
		},
		{
			name: "Solid",
			config: models.Config{
				Framework: models.FrameworkSolid,
			},
			wantJsx: "preserve",
			wantExtraOptions: map[string]interface{}{
				"jsxImportSource": "solid-js",
			},
		},
		{
			name: "Angular",
			config: models.Config{
				Framework: models.FrameworkAngular,
			},
			wantJsx: "preserve",
			wantExtraOptions: map[string]interface{}{
				"experimentalDecorators":  true,
				"emitDecoratorMetadata": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsConfigs := generators.GenerateTSConfig(tt.config)

			// Check base config
			if len(tsConfigs.Base) == 0 {
				t.Error("Base TSConfig should not be empty")
			}

			// Check app config
			appConfig, ok := tsConfigs.App["compilerOptions"].(map[string]interface{})
			if !ok {
				t.Fatal("App TSConfig should have compilerOptions")
			}

			// Verify jsx mode
			if jsx, ok := appConfig["jsx"].(string); !ok || jsx != tt.wantJsx {
				t.Errorf("Expected jsx='%s', got '%v'", tt.wantJsx, appConfig["jsx"])
			}

			// Verify extra options
			for key, wantValue := range tt.wantExtraOptions {
				if gotValue, ok := appConfig[key]; !ok {
					t.Errorf("Missing compiler option '%s'", key)
				} else if gotValue != wantValue {
					t.Errorf("Compiler option '%s': expected %v, got %v", key, wantValue, gotValue)
				}
			}

			// Check node config
			if len(tsConfigs.Node) == 0 {
				t.Error("Node TSConfig should not be empty")
			}
		})
	}
}

func TestGenerateProjectStructure(t *testing.T) {
	tests := []struct {
		name           string
		config         models.Config
		wantDirs       []string
		wantFiles      []string
	}{
		{
			name: "Feature-based",
			config: models.Config{
				Language:  models.LangTypeScript,
				Structure: models.StructureFeatureBased,
			},
			wantDirs: []string{
				"src",
				"public",
				"src/features",
				"src/features/auth",
				"src/features/dashboard",
				"src/components",
				"src/lib",
				"src/hooks",
			},
			wantFiles: []string{
				"src/lib/utils.ts",
				"src/features/README.md",
			},
		},
		{
			name: "Layer-based",
			config: models.Config{
				Language:  models.LangTypeScript,
				Structure: models.StructureLayerBased,
			},
			wantDirs: []string{
				"src",
				"public",
				"src/components",
				"src/pages",
				"src/services",
				"src/utils",
				"src/hooks",
				"src/types",
				"src/lib",
			},
			wantFiles: []string{
				"src/lib/utils.ts",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := testutil.TempDir(t)

			err := generators.GenerateProjectStructure(tempDir, tt.config)
			testutil.AssertNoError(t, err)

			// Verify directories exist
			for _, dir := range tt.wantDirs {
				path := filepath.Join(tempDir, dir)
				testutil.AssertFileExists(t, path)
			}

			// Verify files exist
			for _, file := range tt.wantFiles {
				path := filepath.Join(tempDir, file)
				testutil.AssertFileExists(t, path)
			}
		})
	}
}

func TestGenerateProjectStructure_UtilsFunction(t *testing.T) {
	tests := []struct {
		name     string
		language string
		wantFile string
		wantCode string
	}{
		{
			name:     "TypeScript",
			language: models.LangTypeScript,
			wantFile: "src/lib/utils.ts",
			wantCode: ": (string | undefined)[]",
		},
		{
			name:     "JavaScript",
			language: models.LangJavaScript,
			wantFile: "src/lib/utils.js",
			wantCode: "export function cn(",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := testutil.TempDir(t)

			config := models.Config{
				Language:  tt.language,
				Structure: models.StructureLayerBased,
			}

			err := generators.GenerateProjectStructure(tempDir, config)
			testutil.AssertNoError(t, err)

			utilsPath := filepath.Join(tempDir, tt.wantFile)
			testutil.AssertFileExists(t, utilsPath)
			testutil.AssertFileContains(t, utilsPath, tt.wantCode)
		})
	}
}
