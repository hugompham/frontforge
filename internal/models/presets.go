package models

// QuickPreset returns the opinionated default configuration
func QuickPreset() Config {
	return Config{
		Language:        LangTypeScript,
		Framework:       FrameworkReact,
		PackageManager:  PackageManagerNpm,
		Styling:         StylingTailwind,
		UILibrary:       UILibraryShadcn,
		Routing:         RoutingReactRouter,
		Testing:         TestingVitest,
		StateManagement: StateZustand,
		FormManagement:  FormReactHookForm,
		DataFetching:    DataTanStackQuery,
		Animation:       AnimationFramerMotion,
		Icons:           IconsHeroicons,
		DataViz:         DataVizNone,
		Utilities:       UtilsDateFns,
		I18n:            I18nNone,
		Structure:       StructureFeatureBased,
	}
}
