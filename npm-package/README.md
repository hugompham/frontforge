# frontforge

> A beautiful, interactive CLI for scaffolding modern frontend projects

Also available as `create-frontend-app` for familiar naming.

Built with Go, Bubbletea, Bubbles, and Lipgloss for a polished terminal experience.

## Features

- **Beautiful TUI** - Interactive terminal interface with smooth navigation
- **Fast & Lightweight** - Single native binary, no Node.js runtime needed
- **Modern Frameworks** - React, Vue, Angular, Svelte, Solid, Vanilla
- **Smart Defaults** - Quick mode with opinionated setup
- **Full Control** - Custom mode with 9+ configuration options
- **Latest Packages** - Always uses the newest stable versions
- **Cross-platform** - Works on macOS, Linux, and Windows

## Quick Start

### Using npx (Recommended)

```bash
npx frontforge
# or
npx create-frontend-app
```

### Using npm

```bash
# Install globally
npm install -g frontforge

# Run
frontforge
# or
create-frontend-app
```

### Using yarn

```bash
yarn create frontend-app
```

### Using pnpm

```bash
pnpm create frontend-app
```

## What You Get

The CLI will guide you through an interactive setup to create a fully configured project:

### Quick Mode (Opinionated Defaults)
- TypeScript
- React with Vite
- Tailwind CSS
- React Router
- Vitest for testing
- Zustand for state management
- TanStack Query for data fetching
- Feature-based project structure

### Custom Mode (Full Control)

Choose from:

**Languages:** TypeScript, JavaScript

**Frameworks:** React, Vue 3, Angular, Svelte 5, Solid, Vanilla

**Styling:** Tailwind CSS 4, CSS Modules, Sass, Styled Components, Vanilla CSS

**Routing:** React Router 7, TanStack Router, Vue Router, Angular Router, SvelteKit, Solid Router

**Testing:** Vitest 4, Jest 30, None

**State Management:** Zustand 5, Redux Toolkit, Context API, Pinia 3, Vuex, Svelte Stores, Solid Stores, NgRx

**Data Fetching:** TanStack Query, Axios, Fetch API, SWR, None

**Project Structure:** Feature-based, Layer-based

## Package Versions (December 2025)

All packages use the **latest stable versions**:

- React 19.2.1
- Vue 3.5.13
- Svelte 5.30.0
- Vite 7.2.7
- TypeScript 5.9.3
- Tailwind CSS 4.1.17
- And many more...

See [PACKAGE_VERSIONS.md](https://github.com/hugompham/frontforge/blob/main/PACKAGE_VERSIONS.md) for the complete list.

## Generated Project

After running the CLI, you'll have a complete project with:

- `package.json` with all dependencies
- Build tool configuration (Vite, etc.)
- TypeScript configuration (if selected)
- Tailwind/PostCSS config (if selected)
- Test setup (Vitest/Jest)
- ESLint configuration
- Project directory structure
- Example components and routing
- `.gitignore` and `README.md`

Simply run:

```bash
cd your-project-name
npm install
npm run dev
```

## Requirements

- **Node.js** 18+ (for npm/package management only)
- **macOS**, **Linux**, or **Windows**

Note: The CLI itself is a native Go binary, so you don't need Go installed.

## Supported Platforms

- macOS (Intel and Apple Silicon)
- Linux (x64)
- Windows (x64)

## Why This Tool?

### vs. create-react-app
- Uses Vite instead of Webpack (faster)
- More framework options
- Better TUI experience
- Always up-to-date packages

### vs. create-vite
- More comprehensive setup
- Includes routing, state, testing out of the box
- Beautiful interactive interface
- Feature-based structure option

### vs. Manual Setup
- Saves hours of configuration
- Best practices built-in
- No dependency conflicts
- Tested combinations

## Development

The CLI is built with:

- **Go** - Fast, compiled language
- **Bubbletea** - Elm-inspired TUI framework
- **Bubbles** - Reusable TUI components
- **Lipgloss** - Terminal styling library

This combination provides:
- Native performance
- Beautiful, responsive UI
- Single binary distribution
- No runtime dependencies

## Troubleshooting

### Binary not found after installation

The install script should automatically download the correct binary. If it fails:

1. Check your internet connection
2. Verify your platform is supported
3. Download the binary manually from [GitHub Releases](https://github.com/hugompham/frontforge/releases)

### Permission denied (macOS/Linux)

```bash
chmod +x $(which frontforge)
```

### Command not found

Make sure npm's global bin directory is in your PATH:

```bash
npm config get prefix
```

Add `<prefix>/bin` to your PATH if needed.

## Contributing

Contributions are welcome! Please open an issue or PR on [GitHub](https://github.com/hugompham/frontforge).

## License

MIT

## Links

- [Documentation](https://github.com/hugompham/frontforge)
- [Package Versions](https://github.com/hugompham/frontforge/blob/main/PACKAGE_VERSIONS.md)
- [Changelog](https://github.com/hugompham/frontforge/releases)
- [Issues](https://github.com/hugompham/frontforge/issues)
