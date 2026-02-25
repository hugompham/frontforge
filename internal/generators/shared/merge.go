package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MergePackageJSON reads existing package.json at dir, merges new deps/devDeps/scripts
// without overwriting existing entries. Writes back.
func MergePackageJSON(dir string, deps, devDeps, scripts map[string]string) error {
	pkgPath := filepath.Join(dir, "package.json")

	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Merge dependencies
	if len(deps) > 0 {
		existing := getOrCreateMap(pkg, "dependencies")
		for k, v := range deps {
			if _, exists := existing[k]; !exists {
				existing[k] = v
			}
		}
		pkg["dependencies"] = existing
	}

	// Merge devDependencies
	if len(devDeps) > 0 {
		existing := getOrCreateMap(pkg, "devDependencies")
		for k, v := range devDeps {
			if _, exists := existing[k]; !exists {
				existing[k] = v
			}
		}
		pkg["devDependencies"] = existing
	}

	// Merge scripts
	if len(scripts) > 0 {
		existing := getOrCreateMap(pkg, "scripts")
		for k, v := range scripts {
			if _, exists := existing[k]; !exists {
				existing[k] = v
			}
		}
		pkg["scripts"] = existing
	}

	out, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package.json: %w", err)
	}

	return os.WriteFile(pkgPath, out, 0644)
}

// AddNpmScripts adds scripts to existing package.json (won't overwrite existing keys).
func AddNpmScripts(dir string, scripts map[string]string) error {
	return MergePackageJSON(dir, nil, nil, scripts)
}

func getOrCreateMap(pkg map[string]interface{}, key string) map[string]interface{} {
	if existing, ok := pkg[key].(map[string]interface{}); ok {
		return existing
	}
	return make(map[string]interface{})
}
