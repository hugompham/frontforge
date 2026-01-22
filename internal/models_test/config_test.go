package models_test

import (
	"frontforge/internal/models"
	"testing"
)

func TestFrameworkConstants(t *testing.T) {
	// Verify framework constants are defined
	frameworks := []string{
		models.FrameworkReact,
		models.FrameworkVue,
		models.FrameworkAngular,
		models.FrameworkSvelte,
		models.FrameworkSolid,
		models.FrameworkVanilla,
	}

	for _, fw := range frameworks {
		if fw == "" {
			t.Errorf("framework constant should not be empty")
		}
	}
}

func TestLanguageConstants(t *testing.T) {
	// Verify language constants are defined
	languages := []string{
		models.LangJavaScript,
		models.LangTypeScript,
	}

	for _, lang := range languages {
		if lang == "" {
			t.Errorf("language constant should not be empty")
		}
	}
}

func TestPackageManagerConstants(t *testing.T) {
	// Verify package manager constants are defined
	packageManagers := []string{
		models.PackageManagerNpm,
		models.PackageManagerYarn,
		models.PackageManagerPnpm,
		models.PackageManagerBun,
	}

	for _, pm := range packageManagers {
		if pm == "" {
			t.Errorf("package manager constant should not be empty")
		}
	}
}

func TestStylingConstants(t *testing.T) {
	// Verify styling constants are defined
	stylings := []string{
		models.StylingTailwind,
		models.StylingBootstrap,
		models.StylingCSSModules,
		models.StylingSass,
		models.StylingStyled,
		models.StylingVanilla,
	}

	for _, styling := range stylings {
		if styling == "" {
			t.Errorf("styling constant should not be empty")
		}
	}
}

func TestRoutingConstants(t *testing.T) {
	// Verify routing constants are defined
	routings := []string{
		models.RoutingNone,
		models.RoutingReactRouter,
		models.RoutingVueRouter,
		models.RoutingAngularRouter,
		models.RoutingTanStackRouter,
		models.RoutingFileBased,
		models.RoutingSvelteKit,
		models.RoutingSolidRouter,
	}

	for _, routing := range routings {
		if routing == "" {
			t.Errorf("routing constant should not be empty")
		}
	}
}

func TestTestingConstants(t *testing.T) {
	// Verify testing constants are defined
	testings := []string{
		models.TestingNone,
		models.TestingVitest,
		models.TestingJest,
	}

	for _, testing := range testings {
		if testing == "" {
			t.Errorf("testing constant should not be empty")
		}
	}
}

func TestStateManagementConstants(t *testing.T) {
	// Verify state management constants are defined
	stateManagers := []string{
		models.StateNone,
		models.StateZustand,
		models.StateReduxToolkit,
		models.StateContextAPI,
		models.StatePinia,
		models.StateVuex,
		models.StateSvelteStores,
		models.StateSolidStores,
		models.StateNgRx,
	}

	for _, sm := range stateManagers {
		if sm == "" {
			t.Errorf("state management constant should not be empty")
		}
	}
}

func TestDataFetchingConstants(t *testing.T) {
	// Verify data fetching constants are defined
	dataFetchers := []string{
		models.DataNone,
		models.DataTanStackQuery,
		models.DataFetchAPI,
		models.DataAxios,
		models.DataSWR,
	}

	for _, df := range dataFetchers {
		if df == "" {
			t.Errorf("data fetching constant should not be empty")
		}
	}
}

func TestStructureConstants(t *testing.T) {
	// Verify structure constants are defined
	structures := []string{
		models.StructureFeatureBased,
		models.StructureLayerBased,
	}

	for _, structure := range structures {
		if structure == "" {
			t.Errorf("structure constant should not be empty")
		}
	}
}

func TestConfigStruct(t *testing.T) {
	// Create a sample config
	config := models.Config{
		ProjectName:     "test-project",
		ProjectPath:     "/tmp/test-project",
		Framework:       models.FrameworkReact,
		Language:        models.LangTypeScript,
		PackageManager:  models.PackageManagerNpm,
		Styling:         models.StylingTailwind,
		Routing:         models.RoutingReactRouter,
		Testing:         models.TestingVitest,
		StateManagement: models.StateZustand,
		DataFetching:    models.DataTanStackQuery,
		Structure:       models.StructureFeatureBased,
	}

	// Verify all fields are set
	if config.ProjectName != "test-project" {
		t.Errorf("expected ProjectName=test-project, got %s", config.ProjectName)
	}
	if config.Framework != models.FrameworkReact {
		t.Errorf("expected Framework=react, got %s", config.Framework)
	}
	if config.Language != models.LangTypeScript {
		t.Errorf("expected Language=typescript, got %s", config.Language)
	}
}
