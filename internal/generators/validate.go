package generators

import (
	"encoding/json"
	"fmt"
	"frontforge/internal/models"
	"os"
	"path/filepath"
	"strings"
)

// ValidationResult represents the result of a single validation check
type ValidationResult struct {
	Check   string // Name of the check
	Passed  bool
	Message string // Details about the failure
}

// ValidateProject performs post-generation validation checks
func ValidateProject(projectPath string, config models.Config) []ValidationResult {
	var results []ValidationResult

	// Skip validation in dry-run mode
	if config.DryRun {
		return results
	}

	// 1. Check all expected core files exist and are non-empty
	results = append(results, validateCoreFiles(projectPath, config)...)

	// 2. Validate package.json structure
	results = append(results, validatePackageJSON(projectPath, config)...)

	// 3. Validate index.html references correct main file
	results = append(results, validateIndexHTML(projectPath, config)...)

	// 4. Validate Vite config exists
	results = append(results, validateViteConfig(projectPath, config)...)

	// 5. Validate TypeScript configs when language is TypeScript
	if config.Language == models.LangTypeScript {
		results = append(results, validateTypeScriptConfigs(projectPath)...)
	}

	return results
}

// validateCoreFiles checks that all expected files exist and are non-empty
func validateCoreFiles(projectPath string, config models.Config) []ValidationResult {
	var results []ValidationResult

	expectedFiles := getCoreFiles(config)

	for _, relPath := range expectedFiles {
		fullPath := filepath.Join(projectPath, relPath)
		check := fmt.Sprintf("File exists: %s", relPath)

		info, err := os.Stat(fullPath)
		if err != nil {
			results = append(results, ValidationResult{
				Check:   check,
				Passed:  false,
				Message: "file does not exist",
			})
			continue
		}

		if info.IsDir() {
			results = append(results, ValidationResult{
				Check:   check,
				Passed:  false,
				Message: "expected file, found directory",
			})
			continue
		}

		if info.Size() == 0 {
			results = append(results, ValidationResult{
				Check:   check,
				Passed:  false,
				Message: "file is empty",
			})
			continue
		}

		results = append(results, ValidationResult{
			Check:  check,
			Passed: true,
		})
	}

	return results
}

// validatePackageJSON checks that package.json is valid and has required fields
func validatePackageJSON(projectPath string, config models.Config) []ValidationResult {
	var results []ValidationResult
	packageJSONPath := filepath.Join(projectPath, "package.json")

	// Read and parse package.json
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		results = append(results, ValidationResult{
			Check:   "package.json: readable",
			Passed:  false,
			Message: err.Error(),
		})
		return results
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		results = append(results, ValidationResult{
			Check:   "package.json: valid JSON",
			Passed:  false,
			Message: err.Error(),
		})
		return results
	}

	results = append(results, ValidationResult{
		Check:  "package.json: valid JSON",
		Passed: true,
	})

	// Check required fields
	requiredFields := []string{"name", "scripts", "dependencies"}
	for _, field := range requiredFields {
		check := fmt.Sprintf("package.json: has '%s' field", field)
		if _, exists := pkg[field]; !exists {
			results = append(results, ValidationResult{
				Check:   check,
				Passed:  false,
				Message: fmt.Sprintf("missing required field: %s", field),
			})
		} else {
			results = append(results, ValidationResult{
				Check:  check,
				Passed: true,
			})
		}
	}

	// Check that scripts contains 'dev' command
	if scripts, ok := pkg["scripts"].(map[string]interface{}); ok {
		if _, hasDevScript := scripts["dev"]; !hasDevScript {
			results = append(results, ValidationResult{
				Check:   "package.json: has 'dev' script",
				Passed:  false,
				Message: "missing 'dev' script in package.json",
			})
		} else {
			results = append(results, ValidationResult{
				Check:  "package.json: has 'dev' script",
				Passed: true,
			})
		}
	}

	return results
}

// validateIndexHTML checks that index.html references the correct main file
func validateIndexHTML(projectPath string, config models.Config) []ValidationResult {
	var results []ValidationResult
	indexHTMLPath := filepath.Join(projectPath, "index.html")

	data, err := os.ReadFile(indexHTMLPath)
	if err != nil {
		results = append(results, ValidationResult{
			Check:   "index.html: readable",
			Passed:  false,
			Message: err.Error(),
		})
		return results
	}

	content := string(data)

	// Determine expected main file extension
	ext := getMainFileExtension(config)
	expectedScript := fmt.Sprintf("/src/main.%s", ext)

	if !strings.Contains(content, expectedScript) {
		results = append(results, ValidationResult{
			Check:   "index.html: references main file",
			Passed:  false,
			Message: fmt.Sprintf("does not reference %s", expectedScript),
		})
	} else {
		results = append(results, ValidationResult{
			Check:  "index.html: references main file",
			Passed: true,
		})
	}

	// Check that the referenced main file actually exists
	mainFilePath := filepath.Join(projectPath, "src", fmt.Sprintf("main.%s", ext))
	if _, err := os.Stat(mainFilePath); err != nil {
		results = append(results, ValidationResult{
			Check:   "Main file exists",
			Passed:  false,
			Message: fmt.Sprintf("main.%s does not exist", ext),
		})
	} else {
		results = append(results, ValidationResult{
			Check:  "Main file exists",
			Passed: true,
		})
	}

	return results
}

// validateViteConfig checks that vite.config exists
func validateViteConfig(projectPath string, config models.Config) []ValidationResult {
	var results []ValidationResult

	ext := "js"
	if config.Language == models.LangTypeScript {
		ext = "ts"
	}

	viteConfigPath := filepath.Join(projectPath, fmt.Sprintf("vite.config.%s", ext))

	if _, err := os.Stat(viteConfigPath); err != nil {
		results = append(results, ValidationResult{
			Check:   "vite.config exists",
			Passed:  false,
			Message: fmt.Sprintf("vite.config.%s not found", ext),
		})
	} else {
		results = append(results, ValidationResult{
			Check:  "vite.config exists",
			Passed: true,
		})
	}

	return results
}

// validateTypeScriptConfigs checks that all TypeScript configs exist
func validateTypeScriptConfigs(projectPath string) []ValidationResult {
	var results []ValidationResult

	configs := []string{"tsconfig.json", "tsconfig.app.json", "tsconfig.node.json"}

	for _, configFile := range configs {
		configPath := filepath.Join(projectPath, configFile)
		check := fmt.Sprintf("TypeScript: %s exists", configFile)

		if _, err := os.Stat(configPath); err != nil {
			results = append(results, ValidationResult{
				Check:   check,
				Passed:  false,
				Message: fmt.Sprintf("%s not found", configFile),
			})
		} else {
			results = append(results, ValidationResult{
				Check:  check,
				Passed: true,
			})
		}
	}

	return results
}

// getCoreFiles returns the list of core files that should exist for a given config
func getCoreFiles(config models.Config) []string {
	files := []string{
		"package.json",
		".gitignore",
		"README.md",
		"eslint.config.js",
	}

	// Vite-based frameworks
	files = append(files, "index.html")

	// Add main file
	ext := getMainFileExtension(config)
	files = append(files, filepath.Join("src", fmt.Sprintf("main.%s", ext)))

	// Add App file
	appExt := getAppFileExtension(config)
	if config.Framework == models.FrameworkAngular {
		files = append(files, filepath.Join("src", "app", fmt.Sprintf("app.component.%s", appExt)))
	} else {
		files = append(files, filepath.Join("src", fmt.Sprintf("App.%s", appExt)))
	}

	return files
}

// getAppFileExtension returns the file extension for App file based on config
func getAppFileExtension(config models.Config) string {
	switch config.Framework {
	case models.FrameworkVue:
		return "vue"
	case models.FrameworkSvelte:
		return "svelte"
	case models.FrameworkAngular:
		return "ts"
	default:
		if config.Language == models.LangTypeScript {
			return "tsx"
		}
		return "js"
	}
}

// isViteFramework checks if the framework uses Vite as the direct build tool
func isViteFramework(framework string) bool {
	return true
}
