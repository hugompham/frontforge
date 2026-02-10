package generators

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// DryRunEntry represents a file or directory that would be created
type DryRunEntry struct {
	Path    string // Relative path from project root
	IsDir   bool
	Size    int // Content size in bytes
	Content string
}

// DryRunManifest collects all files/dirs that would be generated
type DryRunManifest struct {
	Entries     []DryRunEntry
	ProjectPath string
	ProjectName string
}

// NewDryRunManifest creates a new manifest
func NewDryRunManifest(projectPath, projectName string) *DryRunManifest {
	return &DryRunManifest{
		Entries:     make([]DryRunEntry, 0),
		ProjectPath: projectPath,
		ProjectName: projectName,
	}
}

// AddFile records a file that would be created
func (m *DryRunManifest) AddFile(path, content string) {
	relPath, _ := filepath.Rel(m.ProjectPath, path)
	m.Entries = append(m.Entries, DryRunEntry{
		Path:    relPath,
		IsDir:   false,
		Size:    len(content),
		Content: content,
	})
}

// AddDir records a directory that would be created
func (m *DryRunManifest) AddDir(path string) {
	relPath, _ := filepath.Rel(m.ProjectPath, path)
	if relPath == "." {
		return // Skip root
	}
	m.Entries = append(m.Entries, DryRunEntry{
		Path:  relPath,
		IsDir: true,
	})
}

// Print outputs the manifest as a tree
func (m *DryRunManifest) Print() {
	fmt.Println()
	fmt.Println("Dry run - files that would be generated:")
	fmt.Println()

	// Build tree structure
	tree := m.buildTree()
	m.printTree(tree, "", true, m.ProjectName)

	// Count files
	fileCount := 0
	for _, entry := range m.Entries {
		if !entry.IsDir {
			fileCount++
		}
	}

	fmt.Println()
	fmt.Printf("%d files would be created\n", fileCount)
	fmt.Println()
}

type treeNode struct {
	name     string
	isDir    bool
	children map[string]*treeNode
}

func (m *DryRunManifest) buildTree() *treeNode {
	root := &treeNode{
		name:     "",
		isDir:    true,
		children: make(map[string]*treeNode),
	}

	for _, entry := range m.Entries {
		parts := strings.Split(entry.Path, string(filepath.Separator))
		current := root

		for i, part := range parts {
			if part == "" || part == "." {
				continue
			}

			if _, exists := current.children[part]; !exists {
				isDir := i < len(parts)-1 || entry.IsDir
				current.children[part] = &treeNode{
					name:     part,
					isDir:    isDir,
					children: make(map[string]*treeNode),
				}
			}
			current = current.children[part]
		}
	}

	return root
}

func (m *DryRunManifest) printTree(node *treeNode, prefix string, isLast bool, name string) {
	if name != "" {
		connector := "├──"
		if isLast {
			connector = "└──"
		}

		suffix := ""
		if node.isDir {
			suffix = "/"
		}

		fmt.Printf("  %s%s %s%s\n", prefix, connector, name, suffix)

		if isLast {
			prefix += "    "
		} else {
			prefix += "│   "
		}
	}

	// Sort children (dirs first, then files, alphabetically)
	children := make([]string, 0, len(node.children))
	for childName := range node.children {
		children = append(children, childName)
	}
	sort.Slice(children, func(i, j int) bool {
		childI := node.children[children[i]]
		childJ := node.children[children[j]]
		if childI.isDir != childJ.isDir {
			return childI.isDir // Dirs first
		}
		return children[i] < children[j] // Alphabetical
	})

	for i, childName := range children {
		child := node.children[childName]
		m.printTree(child, prefix, i == len(children)-1, childName)
	}
}
