package state

// FormState holds all form field values
// This separates form data from the main Model to improve organization
type FormState struct {
	// Setup configuration
	SetupMode string

	// Project basics
	ProjectName string
	Language    string
	Framework   string

	// Tooling
	PackageManager string
	Styling        string
	UILibrary      string

	// Features
	Routing         string
	Testing         string
	StateManagement string
	FormManagement  string
	DataFetching    string

	// Finishing touches
	Animation string
	Icons     string
	DataViz   string
	Utilities string
	I18n      string
	Structure string
}

// NewFormState creates a FormState with recommended defaults
func NewFormState() FormState {
	return FormState{
		SetupMode:       "custom",
		ProjectName:     "my-app",
		Language:        "TypeScript",
		Framework:       "React",
		PackageManager:  "npm",
		Styling:         "Tailwind CSS",
		UILibrary:       "Shadcn/ui",
		Routing:         "React Router",
		Testing:         "Vitest",
		StateManagement: "Zustand",
		FormManagement:  "React Hook Form",
		DataFetching:    "TanStack Query",
		Animation:       "Framer Motion",
		Icons:           "Heroicons",
		DataViz:         "None",
		Utilities:       "date-fns",
		I18n:            "None",
		Structure:       "Feature-based",
	}
}
