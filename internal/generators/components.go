package generators

import (
	"fmt"
	"frontforge/internal/models"
	"strings"
)

// getFileExtension returns the correct file extension based on framework and language
func getFileExtension(config models.Config) string {
	isTS := config.Language == models.LangTypeScript

	switch config.Framework {
	case models.FrameworkVue:
		return "vue"
	case models.FrameworkSvelte:
		return "svelte"
	case models.FrameworkVanilla, models.FrameworkAngular:
		// Vanilla and Angular use plain JS/TS (no JSX)
		if isTS {
			return "ts"
		}
		return "js"
	default:
		// React, Solid use JSX/TSX
		if isTS {
			return "tsx"
		}
		return "jsx"
	}
}

// GenerateIndexHTML creates the index.html file
func GenerateIndexHTML(config models.Config) string {
	ext := getFileExtension(config)

	// Vue uses #app, all other frameworks use #root
	mountID := "root"
	if config.Framework == models.FrameworkVue {
		mountID = "app"
	}

	return fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>%s</title>
  </head>
  <body>
    <div id="%s"></div>
    <script type="module" src="/src/main.%s"></script>
  </body>
</html>
`, config.ProjectName, mountID, ext)
}

// GenerateMainFile creates the main entry file
func GenerateMainFile(config models.Config) string {
	isTS := config.Language == models.LangTypeScript
	ext := getFileExtension(config)

	if config.Framework == models.FrameworkReact {
		var imports strings.Builder
		imports.WriteString("import { StrictMode } from 'react'\n")
		imports.WriteString("import { createRoot } from 'react-dom/client'\n")
		imports.WriteString(fmt.Sprintf("import App from './App.%s'", ext))

		if config.Styling == models.StylingTailwind {
			imports.WriteString("\nimport './index.css'")
		} else if config.Styling == models.StylingBootstrap {
			imports.WriteString("\nimport 'bootstrap/dist/css/bootstrap.min.css'")
		}

		appWrapper := "<App />"

		// Add routing wrapper
		if config.Routing == models.RoutingReactRouter {
			imports.WriteString("\nimport { BrowserRouter } from 'react-router'")
			appWrapper = `<BrowserRouter>
      <App />
    </BrowserRouter>`
		}

		// Add TanStack Query wrapper
		if config.DataFetching == models.DataTanStackQuery {
			imports.WriteString("\nimport { QueryClient, QueryClientProvider } from '@tanstack/react-query'")
			imports.WriteString("\nimport { ReactQueryDevtools } from '@tanstack/react-query-devtools'")
			imports.WriteString("\n\nconst queryClient = new QueryClient()")

			appWrapper = fmt.Sprintf(`<QueryClientProvider client={queryClient}>
      %s
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>`, appWrapper)
		}

		rootSelector := "document.getElementById('root')"
		if isTS {
			rootSelector += "!"
		}

		return fmt.Sprintf(`%s

createRoot(%s).render(
  <StrictMode>
    %s
  </StrictMode>,
)
`, imports.String(), rootSelector, appWrapper)
	} else if config.Framework == models.FrameworkVue {
		var imports strings.Builder
		imports.WriteString("import { createApp } from 'vue'\n")
		imports.WriteString("import App from './App.vue'")

		if config.Styling == models.StylingTailwind {
			imports.WriteString("\nimport './index.css'")
		} else if config.Styling == models.StylingBootstrap {
			imports.WriteString("\nimport 'bootstrap/dist/css/bootstrap.min.css'")
		} else if config.Styling == models.StylingSass {
			imports.WriteString("\nimport './styles.scss'")
		}

		appCode := "const app = createApp(App)"

		if config.Routing == models.RoutingVueRouter {
			imports.WriteString("\nimport router from './router'")
			appCode += "\napp.use(router)"
		}

		if config.StateManagement == models.StatePinia {
			imports.WriteString("\nimport { createPinia } from 'pinia'")
			appCode += "\napp.use(createPinia())"
		}

		return fmt.Sprintf(`%s

%s
app.mount('#app')
`, imports.String(), appCode)
	}

	// Default for other frameworks
	return fmt.Sprintf(`import App from './App.%s'

const root = document.getElementById('root')
if (root) {
  root.innerHTML = '<div id="app"></div>'
}
`, ext)
}

// GenerateAppFile creates the App component
func GenerateAppFile(config models.Config) string {
	isTS := config.Language == models.LangTypeScript

	if config.Framework == models.FrameworkReact {
		var imports strings.Builder
		imports.WriteString("import { useState } from 'react'")

		// Add styling imports
		if config.Styling == models.StylingCSSModules {
			imports.WriteString("\nimport styles from './App.module.css'")
		} else if config.Styling == models.StylingSass {
			imports.WriteString("\nimport './styles.scss'")
		}

		if config.Routing == models.RoutingReactRouter {
			imports.WriteString("\nimport { Routes, Route, Link } from 'react-router'")
		}

		stateExample := ""
		if config.StateManagement == models.StateZustand {
			typeAnnotation := ""
			stateTypeAnnotation := ""
			if isTS {
				typeAnnotation = "<{ count: number; increment: () => void }>"
				stateTypeAnnotation = ": { count: number }"
			}

			imports.WriteString("\nimport { create } from 'zustand'")
			stateExample = fmt.Sprintf(`
const useStore = create%s((set) => ({
  count: 0,
  increment: () => set((state%s) => ({ count: state.count + 1 })),
}))`, typeAnnotation, stateTypeAnnotation)
		}

		routingExample := ""
		if config.Routing == models.RoutingReactRouter {
			routingExample = fmt.Sprintf(`
function Home() {
  return (
    <div>
      <h1>Home Page</h1>
      <p>Welcome to your new %s app!</p>
    </div>
  )
}

function About() {
  return (
    <div>
      <h1>About Page</h1>
      <p>This is the about page.</p>
    </div>
  )
}`, config.ProjectName)
		}

		var appComponent string
		if config.Routing == models.RoutingReactRouter {
			appComponent = `
function App() {
  return (
    <div className="app">
      <nav>
        <Link to="/">Home</Link> | <Link to="/about">About</Link>
      </nav>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/about" element={<About />} />
      </Routes>
    </div>
  )
}`
		} else {
			zustandHook := ""
			zustandButton := ""
			if config.StateManagement == models.StateZustand {
				zustandHook = `
  const { count: zustandCount, increment } = useStore()`
				zustandButton = fmt.Sprintf(`
          <button
            onClick={increment}
            className="%s"
          >
            Zustand Count is {zustandCount}
          </button>`,
					getTailwindButtonClass(config, "green"))
			}

			appComponent = fmt.Sprintf(`
function App() {
  const [count, setCount] = useState(0)%s

  return (
    <div className="%s">
      <div className="%s">
        <h1 className="%s">
          Welcome to %s
        </h1>
        <div className="%s">
          <button
            onClick={() => setCount(count + 1)}
            className="%s"
          >
            Count is {count}
          </button>%s
        </div>
      </div>
    </div>
  )
}`,
				zustandHook,
				getContainerClass(config),
				getCenterClass(config),
				getTitleClass(config),
				config.ProjectName,
				getSpaceClass(config),
				getTailwindButtonClass(config, "blue"),
				zustandButton)
		}

		return fmt.Sprintf(`%s
%s
%s
%s

export default App
`, imports.String(), stateExample, routingExample, appComponent)
	} else if config.Framework == models.FrameworkVue {
		langAttr := ""
		if isTS {
			langAttr = ` lang="ts"`
		}

		return fmt.Sprintf(`<script setup%s>
import { ref } from 'vue'

const count = ref(0)
</script>

<template>
  <div class="%s">
    <div class="%s">
      <h1 class="%s">
        Welcome to %s
      </h1>
      <button
        @click="count++"
        class="%s"
      >
        Count is {{ count }}
      </button>
    </div>
  </div>
</template>

<style scoped>
%s
</style>
`,
			langAttr,
			getContainerClass(config),
			getCenterClass(config),
			getTitleClass(config),
			config.ProjectName,
			getTailwindButtonClass(config, "blue"),
			getVueStyles(config))
	}

	// Default for other frameworks
	return fmt.Sprintf(`function App() {
  return (
    <div>
      <h1>Welcome to %s</h1>
    </div>
  )
}

export default App
`, config.ProjectName)
}

// Helper functions for CSS classes
func getContainerClass(config models.Config) string {
	if config.Styling == models.StylingTailwind {
		return "min-h-screen flex items-center justify-center bg-gray-100"
	} else if config.Styling == models.StylingCSSModules {
		return "styles.app"
	} else if config.Styling == models.StylingSass {
		return "app"
	}
	return "app"
}

func getCenterClass(config models.Config) string {
	if config.Styling == models.StylingTailwind {
		return "text-center"
	}
	return ""
}

func getTitleClass(config models.Config) string {
	if config.Styling == models.StylingTailwind {
		return "text-4xl font-bold mb-4"
	} else if config.Styling == models.StylingCSSModules {
		return "styles.title"
	} else if config.Styling == models.StylingSass {
		return "title"
	}
	return ""
}

func getSpaceClass(config models.Config) string {
	if config.Styling == models.StylingTailwind {
		return "space-y-4"
	}
	return ""
}

func getTailwindButtonClass(config models.Config, color string) string {
	if config.Styling == models.StylingTailwind {
		return fmt.Sprintf("px-4 py-2 bg-%s-500 text-white rounded hover:bg-%s-600", color, color)
	} else if config.Styling == models.StylingCSSModules {
		return "styles.button"
	} else if config.Styling == models.StylingSass {
		return "button"
	}
	return ""
}

func getVueStyles(config models.Config) string {
	if config.Styling == models.StylingVanilla {
		return `.app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  padding: 60px;
}`
	}
	return ""
}
