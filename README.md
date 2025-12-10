# frontforge

A beautiful, interactive CLI for scaffolding modern frontend projects with the latest stable packages.

Built with Go and Charm libraries (Bubbletea, Huh, Lipgloss) for a polished terminal experience.

## Features

- Beautiful terminal UI with smooth animations
- Fast native binary (no Node.js runtime required)
- Smart defaults with complete customization
- Latest stable packages (React 19, Vite 7, TypeScript 5.9, etc.)
- Cross-platform support (macOS, Linux, Windows)
- Quick and Custom setup modes

## Installation

```bash
npx frontforge
```

Or install globally:

```bash
npm install -g frontforge
```

## Usage

```bash
# Interactive mode (creates a new folder with your project name)
frontforge

# Create project in current directory
frontforge -path .

# Create project in specific folder
frontforge -path my-app

# Show help
frontforge -help
```

## Supported Frameworks

- React 19 (with Vite)
- Vue 3.5
- Angular (latest)
- Svelte 5
- Solid
- Vanilla JavaScript/TypeScript

## Package Versions

All package versions are the latest stable releases. See [PACKAGE_VERSIONS.md](./PACKAGE_VERSIONS.md) for the complete list with version numbers and compatibility notes.

## Development

This is the source repository for the frontforge CLI tool.

### Prerequisites

- Go 1.21 or higher
- Node.js 20+ (for npm package publishing)

### Project Structure

```
.
├── internal/
│   ├── generators/     # Project file generators
│   ├── models/         # Data models and constants
│   └── tui/           # Terminal UI (Bubbletea/Huh)
├── npm-package/       # npm wrapper package
├── main.go           # Entry point
├── Makefile          # Build commands
└── PACKAGE_VERSIONS.md # Complete package version list
```

### Building from Source

```bash
# Build for current platform
go build -o frontforge

# Build for all platforms
make build-all

# Run locally
go run main.go
```

### Tech Stack

- **Go** - Fast, compiled language for CLI
- **Bubbletea** - Elm-inspired TUI framework
- **Huh** - Beautiful interactive forms
- **Lipgloss** - Terminal styling and layout
- **Harmonica** - Spring physics for smooth animations

### Publishing

For maintainers only:

1. Update package versions in `internal/generators/package.go`
2. Build binaries: `make build-all`
3. Create GitHub release with binaries
4. Update version in `npm-package/package.json`
5. Publish to npm: `cd npm-package && npm publish`

Note: npm publishing is restricted to authorized maintainers only.

## Contributing

Contributions are welcome. Please ensure:

- Code follows existing patterns and style
- All builds complete successfully (`make build-all`)
- Package versions in `internal/generators/package.go` are up to date

## License

MIT
