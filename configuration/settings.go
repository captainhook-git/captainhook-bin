package configuration

type AppSettings struct {
	AllowFailure     bool
	AnsiColors       bool
	Custom           map[string]string
	FailOnFirstError bool
	GitDirectory     string
	Includes         []string
	IncludeLevel     int
	RunPath          string
	RunAsync         bool
	Verbosity        string
}

func NewDefaultAppSettings() *AppSettings {
	return &AppSettings{
		AllowFailure:     false,
		AnsiColors:       true,
		Custom:           map[string]string{},
		FailOnFirstError: true,
		GitDirectory:     ".git",
		Includes:         []string{},
		IncludeLevel:     1,
		RunAsync:         false,
		Verbosity:        "normal",
	}
}

type ActionSettings struct {
	AllowFailure bool
	WorkingDir   string
	Label        string
}

func NewDefaultActionSettings() *ActionSettings {
	return &ActionSettings{
		AllowFailure: false,
		WorkingDir:   "",
		Label:        "",
	}
}

// NullableAppSettings represents the command argument and options that can be provided or can be skipped
type NullableAppSettings struct {
	AllowFailure     *bool
	AnsiColors       *bool
	Custom           *map[string]string
	FailOnFirstError *bool
	GitDirectory     *string
	Includes         *[]string
	IncludeLevel     *int
	RunPath          *string
	RunAsync         *bool
	Verbosity        *string
}

func NewNullableAppSettings() *NullableAppSettings {
	return &NullableAppSettings{}
}
