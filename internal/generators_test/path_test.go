package generators_test

import (
	"frontforge/internal/generators"
	"frontforge/internal/testutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNormalizePath(t *testing.T) {
	cwd, _ := os.Getwd()

	tests := []struct {
		name      string
		userPath  string
		shouldErr bool
	}{
		{"absolute path", "/tmp/test", false},
		{"relative path", "test", false},
		{"dot dot path", "../test", false},
		{"empty path", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generators.NormalizePath(tt.userPath, cwd)

			if tt.shouldErr {
				testutil.AssertError(t, err)
			} else {
				testutil.AssertNoError(t, err)
				if !filepath.IsAbs(result) {
					t.Errorf("expected absolute path, got %s", result)
				}
			}
		})
	}
}

func TestIsPathSafe(t *testing.T) {
	tests := []struct {
		name string
		path string
		safe bool
	}{
		{"home directory", os.Getenv("HOME"), true},
		{"tmp directory", "/tmp/test", true},
	}

	// Add OS-specific system paths
	if runtime.GOOS != "windows" {
		tests = append(tests, []struct {
			name string
			path string
			safe bool
		}{
			{"root directory", "/", false},
			{"bin directory", "/bin", false},
			{"etc directory", "/etc", false},
			{"usr directory", "/usr", false},
		}...)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generators.IsPathSafe(tt.path)
			if result != tt.safe {
				t.Errorf("IsPathSafe(%q) = %v, want %v", tt.path, result, tt.safe)
			}
		})
	}
}

func TestValidateProjectPath_Empty(t *testing.T) {
	tempDir := testutil.TempDir(t)
	nonExistentPath := filepath.Join(tempDir, "new-project")

	// Non-existent path should be valid
	err := generators.ValidateProjectPath(nonExistentPath)
	testutil.AssertNoError(t, err)
}

func TestValidateProjectPath_Existing(t *testing.T) {
	tempDir := testutil.TempDir(t)

	// Create a file in the directory
	testutil.CreateTempFile(t, tempDir, "test.txt", "content")

	// Directory with files should fail validation
	err := generators.ValidateProjectPath(tempDir)
	testutil.AssertError(t, err)
}

func TestGetProjectName(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/home/user/my-project", "my-project"},
		{"/tmp/test", "test"},
		{"my-project", "my-project"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := generators.GetProjectName(tt.path)
			testutil.AssertEqual(t, result, tt.expected)
		})
	}
}
