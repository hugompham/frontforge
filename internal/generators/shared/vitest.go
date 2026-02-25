package shared

import (
	"fmt"
	"os"
	"path/filepath"
)

// ScaffoldVitest creates vitest.config.ts and test setup files.
// framework should be "nextjs", "sveltekit", or "astro".
func ScaffoldVitest(dir string, framework string) error {
	ext := "ts"

	// Create vitest config
	config := generateVitestConfig(framework)
	if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("vitest.config.%s", ext)), []byte(config), 0644); err != nil {
		return fmt.Errorf("failed to write vitest config: %w", err)
	}

	// Create test directory
	testDir := filepath.Join(dir, "src", "test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return fmt.Errorf("failed to create test directory: %w", err)
	}

	// Create setup file
	setup := generateVitestSetup()
	if err := os.WriteFile(filepath.Join(testDir, fmt.Sprintf("setup.%s", ext)), []byte(setup), 0644); err != nil {
		return fmt.Errorf("failed to write vitest setup: %w", err)
	}

	// Add test deps to package.json
	devDeps := map[string]string{
		"vitest":                    "^4.0.18",
		"@testing-library/jest-dom": "^6.9.1",
		"jsdom":                     "^28.1.0",
	}

	scripts := map[string]string{
		"test": "vitest",
	}

	switch framework {
	case "nextjs":
		devDeps["@testing-library/react"] = "^16.3.2"
		devDeps["@vitejs/plugin-react"] = "^5.1.4"
	case "sveltekit":
		devDeps["@testing-library/svelte"] = "^5.3.1"
	}

	return MergePackageJSON(dir, nil, devDeps, scripts)
}

func generateVitestConfig(framework string) string {
	switch framework {
	case "nextjs":
		return `import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    globals: true,
  },
})
`
	case "sveltekit":
		return `import { defineConfig } from 'vitest/config'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte({ hot: !process.env.VITEST })],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    globals: true,
  },
})
`
	default: // astro
		return `import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    globals: true,
  },
})
`
	}
}

func generateVitestSetup() string {
	return `import { expect } from 'vitest'
import * as matchers from '@testing-library/jest-dom/matchers'

expect.extend(matchers)
`
}
