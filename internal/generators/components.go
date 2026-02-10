package generators

import (
	"fmt"
	"frontforge/internal/models"
	"frontforge/internal/templates"
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

// getMainFileExtension returns the correct file extension for the main entry file
func getMainFileExtension(config models.Config) string {
	isTS := config.Language == models.LangTypeScript

	// Angular always uses .ts (no .tsx)
	if config.Framework == models.FrameworkAngular {
		return "ts"
	}

	// Vanilla uses plain .js or .ts (no JSX)
	if config.Framework == models.FrameworkVanilla {
		if isTS {
			return "ts"
		}
		return "js"
	}

	// All other frameworks (React, Vue, Svelte, Solid) use JSX/TSX for the main file
	if isTS {
		return "tsx"
	}
	return "jsx"
}

// GenerateIndexHTML creates the index.html file
func GenerateIndexHTML(config models.Config) string {
	html, err := templates.Render("static/index.html.tmpl", config)
	if err != nil {
		// Fallback to empty string on error (caller will handle)
		return ""
	}
	return html
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
	} else if config.Framework == models.FrameworkAngular {
		// Angular uses bootstrapApplication with standalone components
		return `import { bootstrapApplication } from '@angular/platform-browser'
import { AppComponent } from './app/app.component'

bootstrapApplication(AppComponent)
  .catch(err => console.error(err))
`
	} else if config.Framework == models.FrameworkSvelte {
		// Svelte 5 uses mount() instead of new App()
		return fmt.Sprintf(`import { mount } from 'svelte'
import App from './App.svelte'

mount(App, {
  target: document.getElementById('root')!
})
`)
	} else if config.Framework == models.FrameworkSolid {
		// Solid uses render from solid-js/web
		return fmt.Sprintf(`import { render } from 'solid-js/web'
import App from './App.%s'

const root = document.getElementById('root')
if (root) {
  render(() => <App />, root)
}
`, ext)
	} else if config.Framework == models.FrameworkVanilla {
		// Vanilla JS/TS - minimal setup
		return fmt.Sprintf(`import App from './App.%s'

const root = document.getElementById('root')
if (root) {
  root.innerHTML = App()
}
`, ext)
	}

	// Fallback (shouldn't reach here)
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
	} else if config.Framework == models.FrameworkAngular {
		// Angular standalone component
		return fmt.Sprintf(`import { Component } from '@angular/core'

@Component({
  selector: 'app-root',
  standalone: true,
  template: `+"`"+`
    <div style="text-align: center; padding: 60px;">
      <h1>Welcome to %s</h1>
      <button (click)="increment()">
        Count is {{ count }}
      </button>
    </div>
  `+"`"+`,
  styles: [`+"`"+`
    button {
      padding: 10px 20px;
      font-size: 16px;
      cursor: pointer;
      background-color: #1976d2;
      color: white;
      border: none;
      border-radius: 4px;
    }
    button:hover {
      background-color: #1565c0;
    }
  `+"`"+`]
})
export class AppComponent {
  count = 0

  increment() {
    this.count++
  }
}
`, config.ProjectName)
	} else if config.Framework == models.FrameworkSvelte {
		// Svelte component (SFC)
		return fmt.Sprintf(`<script>
  let count = $state(0)
</script>

<main>
  <div class="container">
    <h1>Welcome to %s</h1>
    <button onclick={() => count++}>
      Count is {count}
    </button>
  </div>
</main>

<style>
  .container {
    text-align: center;
    padding: 60px;
  }

  h1 {
    font-size: 2rem;
    font-weight: bold;
    margin-bottom: 1rem;
  }

  button {
    padding: 10px 20px;
    font-size: 16px;
    cursor: pointer;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
  }

  button:hover {
    background-color: #0056b3;
  }
</style>
`, config.ProjectName)
	} else if config.Framework == models.FrameworkSolid {
		// Solid component with createSignal
		return fmt.Sprintf(`import { createSignal } from 'solid-js'

function App() {
  const [count, setCount] = createSignal(0)

  return (
    <div style="text-align: center; padding: 60px;">
      <h1>Welcome to %s</h1>
      <button
        onClick={() => setCount(count() + 1)}
        style="padding: 10px 20px; font-size: 16px; cursor: pointer; background-color: #007bff; color: white; border: none; border-radius: 4px;"
      >
        Count is {count()}
      </button>
    </div>
  )
}

export default App
`, config.ProjectName)
	} else if config.Framework == models.FrameworkVanilla {
		// Vanilla JS/TS - returns HTML string
		typeAnnotation := ""
		if isTS {
			typeAnnotation = ": string"
		}
		return fmt.Sprintf(`export default function App()%s {
  return `+"`"+`
    <div style="text-align: center; padding: 60px;">
      <h1>Welcome to %s</h1>
      <button id="counter" style="padding: 10px 20px; font-size: 16px; cursor: pointer; background-color: #007bff; color: white; border: none; border-radius: 4px;">
        Count is 0
      </button>
    </div>
  `+"`"+`
}

// Add event listener after DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  let count = 0
  const button = document.getElementById('counter')
  button?.addEventListener('click', () => {
    count++
    if (button) button.textContent = `+"`Count is ${count}`"+`
  })
})
`, typeAnnotation, config.ProjectName)
	}

	// Fallback (shouldn't reach here)
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
