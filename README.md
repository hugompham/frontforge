# FrontForge

A fast, interactive CLI for scaffolding modern frontend projects with the latest stable packages.

Built with Go and the Charm ecosystem (Bubbletea, Huh, Lipgloss) for a polished terminal experience.

## Features

- Beautiful terminal interface with smooth navigation
- Fast native binary with no Node.js runtime required
- Smart defaults with complete customization options
- Latest stable packages (React 19, Vite 7, TypeScript 5.9, etc.)
- Cross-platform support (macOS, Linux, Windows)
- Quick mode for rapid setup or Custom mode for full control

## Installation

### Using npx (Recommended)

```bash
npx frontforge
```

### Using npm

```bash
# Install globally
npm install -g frontforge

# Run
frontforge
```

### Using yarn

```bash
yarn create frontend-app
```

### Using pnpm

```bash
pnpm create frontend-app
```

## Usage

```bash
# Interactive mode - creates new folder
frontforge

# Create in current directory
frontforge -path .

# Create in specific folder
frontforge -path my-app

# Show help
frontforge -help
```

## What You Get

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

Choose from multiple options for each:

- **Languages**: TypeScript, JavaScript
- **Frameworks**: React, Vue 3, Angular, Svelte 5, Solid, Vanilla, Next.js, Astro, SvelteKit
- **Styling**: Tailwind CSS, CSS Modules, Sass, Styled Components, Vanilla CSS
- **Routing**: React Router, TanStack Router, Vue Router, Angular Router, and more
- **Testing**: Vitest, Jest, or None
- **State Management**: Zustand, Redux Toolkit, Context API, Pinia, Vuex, NgRx, and more
- **Data Fetching**: TanStack Query, Axios, Fetch API, SWR, or None
- **Project Structure**: Feature-based or Layer-based

## Supported Frameworks

### Vite-based
- React 19
- Vue 3.5
- Angular (latest)
- Svelte 5
- Solid
- Vanilla JavaScript/TypeScript

### Meta-frameworks
- Next.js 15 (React, App Router)
- Astro 5 (content-focused)
- SvelteKit 2 (Svelte meta-framework)

## Package Versions

All packages use the latest stable releases. See [PACKAGE_VERSIONS.md](./PACKAGE_VERSIONS.md) for the complete list with version numbers.

## Generated Project Structure

After running FrontForge, you'll have a complete project with:

- `package.json` with all dependencies
- Build tool configuration (Vite, etc.)
- TypeScript configuration (if selected)
- Tailwind/PostCSS config (if selected)
- Test setup (Vitest/Jest)
- ESLint configuration
- Project directory structure
- Example components and routing setup
- `.gitignore` and `README.md`

Simply run:

```bash
cd your-project-name
npm install
npm run dev
```

## Requirements

- Node.js 18+ (for package management only)
- macOS, Linux, or Windows

Note: The CLI itself is a native Go binary, so Go is not required.

## Supported Platforms

- macOS (Intel and Apple Silicon)
- Linux (x64)
- Windows (x64)

## Why FrontForge?

### vs. create-react-app

- Uses Vite instead of Webpack (faster builds)
- More framework options beyond React
- Beautiful terminal interface
- Always up-to-date with latest packages

### vs. create-vite

- More comprehensive out-of-the-box setup
- Includes routing, state management, and testing
- Interactive configuration interface
- Feature-based project structure option

### vs. Manual Setup

- Saves hours of configuration time
- Best practices built-in
- No dependency conflicts
- Tested package combinations

## Roadmap

| Status | Feature | Description |
|--------|---------|-------------|
| Planned | Plugin system | Custom templates/presets via `.frontforge/` config |
| Planned | `frontforge update` | Bump package versions in existing projects |
| Planned | Git repo templates | Scaffold from remote template repos |
| Planned | Monorepo support | Turborepo/Nx workspace scaffolding |
| Planned | Config presets | Save/share custom configurations |
| Planned | Homebrew tap | `brew install frontforge` |
| Planned | VS Code extension | GUI for project scaffolding |
| Planned | `frontforge add` | Add features to existing projects (e.g., add Tailwind) |

## Development

This repository contains the source code for FrontForge.

### Prerequisites

- Go 1.21 or higher
- Node.js 20+ (for npm package publishing)

### Project Structure

```
frontforge/
├── internal/
│   ├── errors/         # Structured error types
│   ├── generators/     # Project file generators
│   ├── logger/         # Structured logging
│   ├── models/         # Data models and constants
│   ├── preflight/      # Pre-flight validation checks
│   └── tui/           # Terminal UI (Bubbletea)
├── npm-package/       # npm wrapper package
├── main.go           # Entry point
└── PACKAGE_VERSIONS.md
```

### Building from Source

```bash
# Build for current platform
go build -o frontforge

# Run locally
go run main.go
```

### Tech Stack

- **Go** - Fast, compiled language for CLI tools
- **Bubbletea** - Elm-inspired terminal UI framework
- **Huh** - Beautiful interactive forms
- **Lipgloss** - Terminal styling and layout

## Contributing

Contributions are welcome! Please ensure:

- Code follows existing patterns and style
- All tests pass (`go test ./...`)
- Builds complete successfully (`go build`)

## Troubleshooting

### Binary not found after installation

If the install script fails to download the binary:

1. Check your internet connection
2. Verify your platform is supported
3. Download manually from [GitHub Releases](https://github.com/hugompham/frontforge/releases)

### Permission denied (macOS/Linux)

```bash
chmod +x $(which frontforge)
```

### Command not found

Ensure npm's global bin directory is in your PATH:

```bash
npm config get prefix
```

Add `<prefix>/bin` to your PATH if needed.

## License

MIT

## Links

- [Package Versions](./PACKAGE_VERSIONS.md)
- [GitHub Issues](https://github.com/hugompham/frontforge/issues)
- [Releases](https://github.com/hugompham/frontforge/releases)
