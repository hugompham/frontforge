package main

import (
	"frontforge/internal/tui"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Define command-line flags
	var projectPath string
	var showHelp bool
	flag.StringVar(&projectPath, "path", "", "Project path (use '.' for current directory, or specify a folder name)")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.BoolVar(&showHelp, "h", false, "Show help information (shorthand)")
	flag.Parse()

	// Show help if requested
	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// Resolve the absolute project path
	var absPath string
	var userPath string
	var err error

	if projectPath == "" {
		// No path specified - will create a new folder based on project name
		// Path will be resolved later after user enters project name in TUI
		absPath = ""
		userPath = ""
	} else if projectPath == "." {
		// Use current directory
		absPath, err = os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		userPath = "."
	} else {
		// Create new directory with specified path
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		absPath = filepath.Join(cwd, projectPath)
		userPath = projectPath
	}

	// Create the Bubbletea program with project path
	p := tea.NewProgram(tui.NewModelWithPath(absPath, userPath))

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

// printHelp displays comprehensive help information
func printHelp() {
	fmt.Println("FRONTFORGE - Modern Frontend Project Scaffolding")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  frontforge [options]")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  -path <path>     Project path (optional)")
	fmt.Println("                   If not specified, creates a new folder with your project name")
	fmt.Println("                   Use '.' for current directory")
	fmt.Println("                   Or specify a folder name for a new directory")
	fmt.Println()
	fmt.Println("  -h, -help        Show this help information")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  Create project (interactive - asks for project name):")
	fmt.Println("    frontforge")
	fmt.Println()
	fmt.Println("  Create project in current directory:")
	fmt.Println("    frontforge -path .")
	fmt.Println()
	fmt.Println("  Create project in specific folder:")
	fmt.Println("    frontforge -path my-app")
	fmt.Println("    frontforge -path my-react-project")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  FrontForge is an interactive CLI tool for scaffolding modern frontend projects.")
	fmt.Println("  It guides you through selecting your framework, tooling, and configuration,")
	fmt.Println("  then generates a production-ready project with best practices.")
	fmt.Println()
	fmt.Println("SUPPORTED FRAMEWORKS:")
	fmt.Println("  - React (with Vite)")
	fmt.Println("  - Vue")
	fmt.Println("  - Angular")
	fmt.Println("  - Svelte")
	fmt.Println("  - Solid")
	fmt.Println("  - Vanilla JavaScript/TypeScript")
	fmt.Println()
	fmt.Println("FEATURES:")
	fmt.Println("  - Interactive TUI with guided configuration")
	fmt.Println("  - TypeScript or JavaScript support")
	fmt.Println("  - Multiple package managers (npm, yarn, pnpm, bun)")
	fmt.Println("  - Popular styling solutions (Tailwind, Bootstrap, Sass, etc.)")
	fmt.Println("  - UI component libraries (Shadcn, MUI, Chakra, Vuetify, etc.)")
	fmt.Println("  - Testing setup (Vitest, Jest)")
	fmt.Println("  - State management (Zustand, Redux, Pinia, etc.)")
	fmt.Println("  - And much more...")
	fmt.Println()
}
