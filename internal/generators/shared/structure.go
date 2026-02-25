package shared

import (
	"os"
	"path/filepath"
)

// ScaffoldFeatureStructure creates feature-based directory layout.
// Adapts to framework (e.g., Next.js uses app/, Astro uses src/pages/).
func ScaffoldFeatureStructure(dir string, framework string) error {
	var dirs []string

	switch framework {
	case "nextjs":
		dirs = []string{
			filepath.Join(dir, "app", "features"),
			filepath.Join(dir, "app", "components"),
			filepath.Join(dir, "lib"),
			filepath.Join(dir, "hooks"),
		}
	case "sveltekit":
		dirs = []string{
			filepath.Join(dir, "src", "lib", "features"),
			filepath.Join(dir, "src", "lib", "components"),
			filepath.Join(dir, "src", "lib", "stores"),
		}
	case "astro":
		dirs = []string{
			filepath.Join(dir, "src", "components"),
			filepath.Join(dir, "src", "layouts"),
			filepath.Join(dir, "src", "styles"),
		}
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	return nil
}
