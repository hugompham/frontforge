package generators

import "frontforge/internal/models"

// PackageJSON represents the structure of package.json
type PackageJSON struct {
	Name            string            `json:"name"`
	Private         bool              `json:"private"`
	Version         string            `json:"version"`
	Type            string            `json:"type"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// GeneratePackageJSON creates a package.json configuration
func GeneratePackageJSON(config models.Config) PackageJSON {
	pkg := PackageJSON{
		Name:            config.ProjectName,
		Private:         true,
		Version:         "0.0.0",
		Type:            "module",
		Scripts:         make(map[string]string),
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	// Scripts based on framework
	if config.Framework == models.FrameworkReact || config.Framework == models.FrameworkVue || config.Framework == models.FrameworkSvelte {
		pkg.Scripts["dev"] = "vite"
		pkg.Scripts["build"] = "vite build"
		pkg.Scripts["preview"] = "vite preview"
	} else if config.Framework == models.FrameworkAngular {
		pkg.Scripts["dev"] = "ng serve"
		pkg.Scripts["build"] = "ng build"
		pkg.Scripts["test"] = "ng test"
	}

	// Lint script
	pkg.Scripts["lint"] = "eslint ."

	// Test script
	if config.Testing == models.TestingVitest {
		pkg.Scripts["test"] = "vitest"
	} else if config.Testing == models.TestingJest {
		pkg.Scripts["test"] = "jest"
	}

	// Framework dependencies
	switch config.Framework {
	case models.FrameworkReact:
		pkg.Dependencies["react"] = "^19.2.1"
		pkg.Dependencies["react-dom"] = "^19.2.1"
		pkg.DevDependencies["@vitejs/plugin-react"] = "^5.1.2"
		pkg.DevDependencies["@types/react"] = "^19.2.7"
		pkg.DevDependencies["@types/react-dom"] = "^19.2.3"
	case models.FrameworkVue:
		pkg.Dependencies["vue"] = "^3.5.13"
		pkg.DevDependencies["@vitejs/plugin-vue"] = "^6.0.2"
	case models.FrameworkSvelte:
		pkg.Dependencies["svelte"] = "^5.30.0"
		pkg.DevDependencies["@sveltejs/vite-plugin-svelte"] = "^6.2.1"
	}

	// Vite
	if config.Framework == models.FrameworkReact || config.Framework == models.FrameworkVue || config.Framework == models.FrameworkSvelte {
		pkg.DevDependencies["vite"] = "^7.2.7"
	}

	// TypeScript
	if config.Language == models.LangTypeScript {
		pkg.DevDependencies["typescript"] = "^5.9.3"
		if config.Framework == models.FrameworkReact {
			pkg.DevDependencies["typescript-eslint"] = "^8.49.0"
		}
	}

	// Routing
	switch config.Routing {
	case models.RoutingReactRouter:
		pkg.Dependencies["react-router"] = "^7.10.1"
	case models.RoutingTanStackRouter:
		pkg.Dependencies["@tanstack/react-router"] = "^1.140.2"
	case models.RoutingVueRouter:
		pkg.Dependencies["vue-router"] = "^4.6.3"
	}

	// Styling
	switch config.Styling {
	case models.StylingTailwind:
		pkg.DevDependencies["tailwindcss"] = "^4.1.18"
		pkg.DevDependencies["@tailwindcss/vite"] = "^4.1.18"
	case models.StylingBootstrap:
		pkg.Dependencies["bootstrap"] = "^5.3.3"
	case models.StylingSass:
		pkg.DevDependencies["sass"] = "^1.95.0"
	case models.StylingStyled:
		pkg.Dependencies["styled-components"] = "^6.1.15"
	}

	// State Management
	switch config.StateManagement {
	case models.StateZustand:
		pkg.Dependencies["zustand"] = "^5.0.9"
	case models.StateReduxToolkit:
		pkg.Dependencies["@reduxjs/toolkit"] = "^2.11.1"
		pkg.Dependencies["react-redux"] = "^9.2.0"
	case models.StatePinia:
		pkg.Dependencies["pinia"] = "^3.0.4"
	}

	// Data Fetching
	switch config.DataFetching {
	case models.DataTanStackQuery:
		pkg.Dependencies["@tanstack/react-query"] = "^5.90.12"
		pkg.DevDependencies["@tanstack/react-query-devtools"] = "^5.91.1"
	case models.DataAxios:
		pkg.Dependencies["axios"] = "^1.7.9"
	case models.DataSWR:
		pkg.Dependencies["swr"] = "^2.3.2"
	}

	// Testing
	if config.Testing == models.TestingVitest {
		pkg.DevDependencies["vitest"] = "^4.0.15"
		pkg.DevDependencies["@testing-library/react"] = "^16.3.0"
		pkg.DevDependencies["@testing-library/jest-dom"] = "^6.9.1"
		pkg.DevDependencies["jsdom"] = "^25.0.1"
	} else if config.Testing == models.TestingJest {
		pkg.DevDependencies["jest"] = "^30.2.0"
		pkg.DevDependencies["@testing-library/react"] = "^16.3.0"
		pkg.DevDependencies["@testing-library/jest-dom"] = "^6.9.1"
	}

	// ESLint
	pkg.DevDependencies["eslint"] = "^9.39.1"
	pkg.DevDependencies["@eslint/js"] = "^9.39.1"
	pkg.DevDependencies["globals"] = "^15.15.0"

	if config.Framework == models.FrameworkReact {
		pkg.DevDependencies["eslint-plugin-react-hooks"] = "^7.0.1"
		pkg.DevDependencies["eslint-plugin-react-refresh"] = "^0.4.24"
	}

	// UI Component Libraries
	switch config.UILibrary {
	case models.UILibraryShadcn:
		// Shadcn requires manual setup, add base dependencies
		pkg.Dependencies["class-variance-authority"] = "^0.7.1"
		pkg.Dependencies["clsx"] = "^2.1.1"
		pkg.Dependencies["tailwind-merge"] = "^3.4.0"
		pkg.Dependencies["@radix-ui/react-slot"] = "^1.1.1"
	case models.UILibraryMUI:
		pkg.Dependencies["@mui/material"] = "^7.4.0"
		pkg.Dependencies["@emotion/react"] = "^11.14.0"
		pkg.Dependencies["@emotion/styled"] = "^11.14.0"
	case models.UILibraryChakra:
		pkg.Dependencies["@chakra-ui/react"] = "^3.4.0"
		pkg.Dependencies["@emotion/react"] = "^11.14.0"
		pkg.Dependencies["@emotion/styled"] = "^11.14.0"
	case models.UILibraryAntD:
		pkg.Dependencies["antd"] = "^6.0.0"
	case models.UILibraryHeadless:
		pkg.Dependencies["@headlessui/react"] = "^2.2.9"
	case models.UILibraryVuetify:
		pkg.Dependencies["vuetify"] = "^3.7.7"
	case models.UILibraryPrimeVue:
		pkg.Dependencies["primevue"] = "^4.3.0"
	case models.UILibraryElementUI:
		pkg.Dependencies["element-plus"] = "^2.9.4"
	case models.UILibraryNaiveUI:
		pkg.Dependencies["naive-ui"] = "^2.41.0"
	case models.UILibraryAngularMaterial:
		pkg.Dependencies["@angular/material"] = "^19.2.0"
	case models.UILibraryPrimeNG:
		pkg.Dependencies["primeng"] = "^21.0.1"
	case models.UILibraryNGZorro:
		pkg.Dependencies["ng-zorro-antd"] = "^19.2.0"
	}

	// Form Management
	switch config.FormManagement {
	case models.FormReactHookForm:
		pkg.Dependencies["react-hook-form"] = "^7.68.0"
		pkg.Dependencies["@hookform/resolvers"] = "^3.11.2"
		pkg.Dependencies["zod"] = "^4.1.13"
	case models.FormFormik:
		pkg.Dependencies["formik"] = "^2.4.9"
		pkg.Dependencies["yup"] = "^1.7.1"
	case models.FormTanStackForm:
		pkg.Dependencies["@tanstack/react-form"] = "^0.42.0"
	case models.FormVeeValidate:
		pkg.Dependencies["vee-validate"] = "^4.15.1"
		pkg.Dependencies["yup"] = "^1.7.1"
	case models.FormZod:
		pkg.Dependencies["zod"] = "^4.1.13"
	case models.FormYup:
		pkg.Dependencies["yup"] = "^1.7.1"
	}

	// Animation
	switch config.Animation {
	case models.AnimationFramerMotion:
		pkg.Dependencies["motion"] = "^12.34.0"
	case models.AnimationGSAP:
		pkg.Dependencies["gsap"] = "^3.14.1"
	case models.AnimationAutoAnimate:
		pkg.Dependencies["@formkit/auto-animate"] = "^0.9.2"
	case models.AnimationReactSpring:
		pkg.Dependencies["@react-spring/web"] = "^9.8.2"
	}

	// Icons
	switch config.Icons {
	case models.IconsReactIcons:
		pkg.Dependencies["react-icons"] = "^5.4.0"
	case models.IconsVueIcons:
		pkg.Dependencies["@vicons/ionicons5"] = "^0.12.0"
	case models.IconsHeroicons:
		pkg.Dependencies["@heroicons/react"] = "^2.2.0"
	case models.IconsLucide:
		pkg.Dependencies["lucide-react"] = "^0.469.0"
	case models.IconsFontAwesome:
		pkg.Dependencies["@fortawesome/fontawesome-svg-core"] = "^6.7.2"
		pkg.Dependencies["@fortawesome/free-solid-svg-icons"] = "^6.7.2"
		pkg.Dependencies["@fortawesome/react-fontawesome"] = "^0.2.3"
	}

	// Data Visualization
	switch config.DataViz {
	case models.DataVizRecharts:
		pkg.Dependencies["recharts"] = "^3.5.1"
	case models.DataVizChartJS:
		pkg.Dependencies["chart.js"] = "^4.5.0"
		pkg.Dependencies["react-chartjs-2"] = "^5.3.0"
	case models.DataVizECharts:
		pkg.Dependencies["echarts"] = "^6.0.0"
		pkg.Dependencies["echarts-for-react"] = "^3.0.2"
	case models.DataVizNivo:
		pkg.Dependencies["@nivo/core"] = "^0.89.0"
		pkg.Dependencies["@nivo/line"] = "^0.89.0"
		pkg.Dependencies["@nivo/bar"] = "^0.89.0"
	}

	// Utilities
	switch config.Utilities {
	case models.UtilsDateFns:
		pkg.Dependencies["date-fns"] = "^4.1.0"
	case models.UtilsDayJS:
		pkg.Dependencies["dayjs"] = "^1.11.14"
	case models.UtilsLodash:
		pkg.Dependencies["lodash-es"] = "^4.17.21"
		pkg.DevDependencies["@types/lodash-es"] = "^4.17.12"
	}

	// Internationalization
	switch config.I18n {
	case models.I18nReactI18next:
		pkg.Dependencies["react-i18next"] = "^16.4.1"
		pkg.Dependencies["i18next"] = "^25.7.2"
	case models.I18nVueI18n:
		pkg.Dependencies["vue-i18n"] = "^10.0.8"
	}

	return pkg
}
