package generators_test

import (
	"frontforge/internal/generators/meta"
	"frontforge/internal/models"
	"testing"

	// Trigger init() registration.
	_ "frontforge/internal/generators/astro"
	_ "frontforge/internal/generators/nextjs"
	_ "frontforge/internal/generators/sveltekit"
)

func TestMetaGeneratorInterfaceCompliance(t *testing.T) {
	t.Parallel()

	frameworks := []struct {
		name               string
		hasStyling         bool
		hasTesting         bool
		hasStateManagement bool
		hasDataFetching    bool
	}{
		{models.FrameworkNextJS, true, true, true, true},
		{models.FrameworkAstro, true, true, false, false},
		{models.FrameworkSvelteKit, true, true, true, true},
	}

	for _, fw := range frameworks {
		t.Run(fw.name, func(t *testing.T) {
			t.Parallel()

			gen, ok := meta.Get(fw.name)
			if !ok {
				t.Fatalf("framework %q not registered; init() did not fire", fw.name)
			}

			opts := gen.SupportedOptions()

			if fw.hasStyling && len(opts.Styling) == 0 {
				t.Errorf("SupportedOptions().Styling is empty for %s", fw.name)
			}
			if fw.hasTesting && len(opts.Testing) == 0 {
				t.Errorf("SupportedOptions().Testing is empty for %s", fw.name)
			}
			if fw.hasStateManagement && len(opts.StateManagement) == 0 {
				t.Errorf("SupportedOptions().StateManagement is empty for %s", fw.name)
			}
			if !fw.hasStateManagement && opts.StateManagement != nil {
				t.Errorf("SupportedOptions().StateManagement should be nil for %s", fw.name)
			}
			if fw.hasDataFetching && len(opts.DataFetching) == 0 {
				t.Errorf("SupportedOptions().DataFetching is empty for %s", fw.name)
			}
			if !fw.hasDataFetching && opts.DataFetching != nil {
				t.Errorf("SupportedOptions().DataFetching should be nil for %s", fw.name)
			}
		})
	}
}
