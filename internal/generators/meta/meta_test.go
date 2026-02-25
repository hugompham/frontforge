package meta

import (
	"frontforge/internal/models"
	"strings"
	"testing"
)

// --- helpers ---

// stubGenerator is a minimal MetaGenerator for testing.
type stubGenerator struct {
	scaffoldErr     error
	postScaffoldErr error
	scaffoldCalled  bool
	postCalled      bool
}

func (s *stubGenerator) Scaffold(cfg models.Config) error {
	s.scaffoldCalled = true
	return s.scaffoldErr
}

func (s *stubGenerator) PostScaffold(cfg models.Config) error {
	s.postCalled = true
	return s.postScaffoldErr
}

func (s *stubGenerator) SupportedOptions() OptionMatrix {
	return OptionMatrix{Styling: []string{"Vanilla CSS"}}
}

func (s *stubGenerator) ProbeVersion() string { return "0.0.0-stub" }

// saveAndRestore snapshots the registry, returning a cleanup func that restores it.
func saveAndRestore() func() {
	snapshot := make(map[string]MetaGenerator, len(generators))
	for k, v := range generators {
		snapshot[k] = v
	}
	return func() {
		generators = snapshot
	}
}

// --- Register / Get ---

func TestRegisterAndGet(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	// Pre-register stubs for table-driven cases.
	stubA := &stubGenerator{}
	stubB := &stubGenerator{}
	Register("FrameworkA", stubA)
	Register("FrameworkB", stubB)

	tests := []struct {
		name      string
		framework string
		wantOK    bool
	}{
		{"registered A", "FrameworkA", true},
		{"registered B", "FrameworkB", true},
		{"unknown framework", "UnknownJS", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, ok := Get(tt.framework)
			if ok != tt.wantOK {
				t.Fatalf("Get(%q) ok = %v, want %v", tt.framework, ok, tt.wantOK)
			}
			if tt.wantOK && gen == nil {
				t.Fatalf("Get(%q) returned nil generator with ok=true", tt.framework)
			}
			if !tt.wantOK && gen != nil {
				t.Fatalf("Get(%q) returned non-nil generator with ok=false", tt.framework)
			}
		})
	}
}

func TestRegisterOverwrite(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	stub := &stubGenerator{}
	Register("TestFramework", stub)

	gen, ok := Get("TestFramework")
	if !ok || gen != stub {
		t.Fatal("expected stub generator after Register")
	}

	stub2 := &stubGenerator{}
	Register("TestFramework", stub2)

	gen, ok = Get("TestFramework")
	if !ok || gen != stub2 {
		t.Fatal("Register should overwrite previous generator")
	}
}

// --- RunMetaScaffold ---

func TestRunMetaScaffold_UnregisteredFramework(t *testing.T) {
	cfg := models.Config{Framework: "NoSuchFramework"}
	err := RunMetaScaffold(cfg)
	if err == nil {
		t.Fatal("expected error for unregistered framework")
	}

	se, ok := err.(*ScaffoldError)
	if !ok {
		t.Fatalf("expected *ScaffoldError, got %T", err)
	}
	if se.Framework != "NoSuchFramework" {
		t.Errorf("ScaffoldError.Framework = %q, want %q", se.Framework, "NoSuchFramework")
	}
	if se.ExitCode != -1 {
		t.Errorf("ScaffoldError.ExitCode = %d, want -1", se.ExitCode)
	}
}

func TestRunMetaScaffold_CallsScaffoldAndPostScaffold(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	stub := &stubGenerator{}
	Register("stub-fw", stub)

	cfg := models.Config{Framework: "stub-fw"}
	err := RunMetaScaffold(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stub.scaffoldCalled {
		t.Error("Scaffold was not called")
	}
	if !stub.postCalled {
		t.Error("PostScaffold was not called")
	}
}

func TestRunMetaScaffold_NoScaffoldFlag(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	stub := &stubGenerator{}
	Register("stub-fw", stub)

	cfg := models.Config{Framework: "stub-fw", NoScaffold: true}
	err := RunMetaScaffold(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.scaffoldCalled {
		t.Error("Scaffold should be skipped when NoScaffold=true")
	}
	if !stub.postCalled {
		t.Error("PostScaffold should still be called when NoScaffold=true")
	}
}

func TestRunMetaScaffold_ScaffoldError(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	stub := &stubGenerator{scaffoldErr: &ScaffoldError{
		Framework: "stub-fw",
		Command:   "npx stub",
		ExitCode:  1,
		Stderr:    "boom",
	}}
	Register("stub-fw", stub)

	cfg := models.Config{Framework: "stub-fw"}
	err := RunMetaScaffold(cfg)
	if err == nil {
		t.Fatal("expected error when Scaffold fails")
	}
	if !stub.scaffoldCalled {
		t.Error("Scaffold should have been called")
	}
	if stub.postCalled {
		t.Error("PostScaffold should not be called after Scaffold error")
	}
}

func TestRunMetaScaffold_PostScaffoldError(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	stub := &stubGenerator{postScaffoldErr: &ScaffoldError{
		Framework: "stub-fw",
		Command:   "npm install",
		ExitCode:  2,
		Stderr:    "install failed",
	}}
	Register("stub-fw", stub)

	cfg := models.Config{Framework: "stub-fw"}
	err := RunMetaScaffold(cfg)
	if err == nil {
		t.Fatal("expected error when PostScaffold fails")
	}
}

// --- ScaffoldError.Error ---

func TestScaffoldError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      ScaffoldError
		contains []string
	}{
		{
			name: "no command (registry miss)",
			err: ScaffoldError{
				Framework: "Next.js",
				Command:   "",
				ExitCode:  -1,
				Stderr:    "no generator registered for framework",
			},
			contains: []string{"scaffold error for Next.js", "no generator registered"},
		},
		{
			name: "with command and exit code",
			err: ScaffoldError{
				Framework: "Astro",
				Command:   "npm create astro@latest",
				ExitCode:  127,
				Stderr:    "command not found",
			},
			contains: []string{"scaffold error for Astro", "exit 127", "command not found", "Command: npm create astro@latest"},
		},
		{
			name: "exit code zero with stderr",
			err: ScaffoldError{
				Framework: "SvelteKit",
				Command:   "npx sv create",
				ExitCode:  0,
				Stderr:    "warning",
			},
			contains: []string{"exit 0", "warning", "Command: npx sv create"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			msg := tt.err.Error()
			for _, substr := range tt.contains {
				if !strings.Contains(msg, substr) {
					t.Errorf("Error() = %q, want substring %q", msg, substr)
				}
			}
		})
	}
}

// --- Interface compliance ---
//
// Real framework generators (nextjs, astro, sveltekit) satisfy MetaGenerator
// at compile time -- their init() calls meta.Register which accepts MetaGenerator.
// We cannot import them here (import cycle), so we verify via a compile-time
// assertion on our stub and a runtime check that registered generators work.

// Compile-time proof that *stubGenerator satisfies MetaGenerator.
var _ MetaGenerator = (*stubGenerator)(nil)

func TestRegisteredGeneratorMethods(t *testing.T) {
	restore := saveAndRestore()
	defer restore()

	stub := &stubGenerator{}
	Register("compliance-check", stub)

	gen, ok := Get("compliance-check")
	if !ok {
		t.Fatal("expected registered generator")
	}

	// SupportedOptions
	opts := gen.SupportedOptions()
	if len(opts.Styling) == 0 {
		t.Error("SupportedOptions().Styling should not be empty")
	}

	// ProbeVersion
	ver := gen.ProbeVersion()
	if ver != "0.0.0-stub" {
		t.Errorf("ProbeVersion() = %q, want %q", ver, "0.0.0-stub")
	}
}
