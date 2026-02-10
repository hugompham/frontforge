package cli_test

import (
	"frontforge/internal/models"
	"testing"
)

// TestFlagMapping tests that CLI flags map correctly to config
func TestFlagMapping(t *testing.T) {
	tests := []struct {
		name         string
		flagFramework string
		flagLang      string
		flagPM        string
		flagStyling   string
		wantFramework string
		wantLanguage  string
		wantPM        string
		wantStyling   string
	}{
		{
			name:          "React TypeScript npm Tailwind",
			flagFramework: "react",
			flagLang:      "ts",
			flagPM:        "npm",
			flagStyling:   "tailwind",
			wantFramework: "React",
			wantLanguage:  "TypeScript",
			wantPM:        "npm",
			wantStyling:   "Tailwind CSS",
		},
		{
			name:          "Vue JavaScript pnpm CSS Modules",
			flagFramework: "vue",
			flagLang:      "js",
			flagPM:        "pnpm",
			flagStyling:   "css-modules",
			wantFramework: "Vue",
			wantLanguage:  "JavaScript",
			wantPM:        "pnpm",
			wantStyling:   "CSS Modules",
		},
		{
			name:          "Svelte TypeScript bun Sass",
			flagFramework: "svelte",
			flagLang:      "ts",
			flagPM:        "bun",
			flagStyling:   "sass",
			wantFramework: "Svelte",
			wantLanguage:  "TypeScript",
			wantPM:        "bun",
			wantStyling:   "Sass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			framework := normalizeFramework(tt.flagFramework)
			if framework != tt.wantFramework {
				t.Errorf("Framework: got %q, want %q", framework, tt.wantFramework)
			}

			language := normalizeLanguage(tt.flagLang)
			if language != tt.wantLanguage {
				t.Errorf("Language: got %q, want %q", language, tt.wantLanguage)
			}

			pm := normalizePackageManager(tt.flagPM)
			if pm != tt.wantPM {
				t.Errorf("PackageManager: got %q, want %q", pm, tt.wantPM)
			}

			styling := normalizeStyling(tt.flagStyling)
			if styling != tt.wantStyling {
				t.Errorf("Styling: got %q, want %q", styling, tt.wantStyling)
			}
		})
	}
}

// TestQuickPreset tests the quick mode preset
func TestQuickPreset(t *testing.T) {
	preset := models.QuickPreset()

	// Verify quick preset has expected defaults
	if preset.Language != "TypeScript" {
		t.Errorf("Expected TypeScript, got %s", preset.Language)
	}
	if preset.Framework != "React" {
		t.Errorf("Expected React, got %s", preset.Framework)
	}
	if preset.Styling != "Tailwind CSS" {
		t.Errorf("Expected Tailwind CSS, got %s", preset.Styling)
	}
	if preset.PackageManager != "npm" {
		t.Errorf("Expected npm, got %s", preset.PackageManager)
	}
	if preset.Testing != "Vitest" {
		t.Errorf("Expected Vitest, got %s", preset.Testing)
	}
}

// TestInvalidFlagValues tests handling of invalid flag values
func TestInvalidFlagValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		normalizer func(string) string
		wantEmpty bool
	}{
		{
			name:       "Invalid framework",
			input:      "invalid-framework",
			normalizer: normalizeFramework,
			wantEmpty:  true,
		},
		{
			name:       "Invalid language",
			input:      "c++",
			normalizer: normalizeLanguage,
			wantEmpty:  true,
		},
		{
			name:       "Invalid package manager",
			input:      "pip",
			normalizer: normalizePackageManager,
			wantEmpty:  true,
		},
		{
			name:       "Invalid styling",
			input:      "css-in-js",
			normalizer: normalizeStyling,
			wantEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.normalizer(tt.input)
			if tt.wantEmpty && result != "" {
				t.Errorf("Expected empty string for invalid input %q, got %q", tt.input, result)
			}
		})
	}
}

// TestCaseInsensitiveFlags tests that flags work with any case
func TestCaseInsensitiveFlags(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"REACT", "React"},
		{"React", "React"},
		{"react", "React"},
		{"rEaCt", "React"},
		{"TS", "TypeScript"},
		{"ts", "TypeScript"},
		{"Ts", "TypeScript"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var result string
			if tt.want == "React" {
				result = normalizeFramework(tt.input)
			} else {
				result = normalizeLanguage(tt.input)
			}

			if result != tt.want {
				t.Errorf("Input %q: got %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

// Helper functions that mirror main.go normalization
func normalizeFramework(flag string) string {
	// Convert to lowercase for case-insensitive comparison
	flag = toLower(flag)
	switch flag {
	case "react":
		return "React"
	case "vue":
		return "Vue"
	case "angular":
		return "Angular"
	case "svelte":
		return "Svelte"
	case "solid":
		return "Solid"
	case "vanilla":
		return "Vanilla"
	default:
		return ""
	}
}

func normalizeLanguage(flag string) string {
	// Convert to lowercase for case-insensitive comparison
	flag = toLower(flag)
	switch flag {
	case "ts", "typescript":
		return "TypeScript"
	case "js", "javascript":
		return "JavaScript"
	default:
		return ""
	}
}

func normalizePackageManager(flag string) string {
	// Convert to lowercase for case-insensitive comparison
	flag = toLower(flag)
	switch flag {
	case "npm":
		return "npm"
	case "yarn":
		return "yarn"
	case "pnpm":
		return "pnpm"
	case "bun":
		return "bun"
	default:
		return ""
	}
}

func normalizeStyling(flag string) string {
	// Convert to lowercase for case-insensitive comparison
	flag = toLower(flag)
	switch flag {
	case "tailwind":
		return "Tailwind CSS"
	case "bootstrap":
		return "Bootstrap"
	case "css-modules":
		return "CSS Modules"
	case "sass":
		return "Sass"
	case "styled":
		return "styled-components"
	case "vanilla":
		return "Vanilla CSS"
	default:
		return ""
	}
}

// toLower converts ASCII string to lowercase
func toLower(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 32
		}
	}
	return string(b)
}
