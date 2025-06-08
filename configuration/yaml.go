package configuration

type YamlConfiguration struct {
	Settings *YamlAppSettings     `yaml:"config,omitempty"`
	Hooks    *map[string]YamlHook `yaml:"hooks"`
}

type YamlAppSettings struct {
	AllowFailure     *bool              `yaml:"allow-failure,omitempty"`
	AnsiColors       *bool              `yaml:"ansi-colors,omitempty"`
	Custom           *map[string]string `yaml:"custom,omitempty"`
	FailOnFirstError *bool              `yaml:"fail-on-first-error,omitempty"`
	GitDirectory     *string            `yaml:"git-directory,omitempty"`
	Includes         *[]string          `yaml:"includes,omitempty"`
	IncludeLevel     *int               `yaml:"includes-level,omitempty"`
	RunAsync         *bool              `yaml:"run-async,omitempty"`
	RunPath          *string            `yaml:"run-path,omitempty"`
	Verbosity        *string            `yaml:"verbosity,omitempty"`
}

type YamlHook struct {
	Actions []YamlAction `yaml:"actions"`
}

type YamlAction struct {
	Run        string                  `yaml:"run"`
	Options    *map[string]interface{} `yaml:"options,omitempty"`
	Settings   *YamlActionSettings     `yaml:"config,omitempty"`
	Conditions []YamlCondition         `yaml:"conditions,omitempty"`
}

type YamlActionSettings struct {
	Label        *string `yaml:"label,omitempty"`
	AllowFailure *bool   `yaml:"failure-allowed,omitempty"`
	RunAsync     *bool   `yaml:"run-async,omitempty"`
	WorkingDir   *string `yaml:"working-directory,omitempty"`
}
type YamlCondition struct {
	Run        string                  `yaml:"run"`
	Options    *map[string]interface{} `yaml:"options,omitempty"`
	Conditions []YamlCondition         `yaml:"conditions,omitempty"`
}

func createActionFromYaml(yaml YamlAction) *Action {
	return &Action{
		run:        yaml.Run,
		settings:   createActionSettingsFromYaml(yaml.Settings),
		conditions: createConditionsFromYaml(yaml.Conditions),
		options:    createOptionsFromYaml(yaml.Options),
	}
}

func createActionSettingsFromYaml(yaml *YamlActionSettings) *ActionSettings {
	a := NewDefaultActionSettings()
	if yaml == nil {
		return a
	}
	if yaml.AllowFailure != nil {
		a.AllowFailure = *yaml.AllowFailure
	}
	if yaml.WorkingDir != nil {
		a.WorkingDir = *yaml.WorkingDir
	}
	if yaml.Label != nil {
		a.Label = *yaml.Label
	}
	return a
}

func createConditionsFromYaml(yamlConditions []YamlCondition) []*Condition {
	var conditions []*Condition
	if yamlConditions == nil {
		return conditions
	}
	for _, condition := range yamlConditions {
		conditions = append(conditions, createConditionFromYaml(condition))
	}
	return conditions
}

func createConditionFromYaml(yaml YamlCondition) *Condition {
	var c []*Condition

	// default value empty options
	opts := map[string]interface{}{}
	o := NewOptions(opts)

	if yaml.Options != nil {
		o = createOptionsFromJson(yaml.Options)
	}
	if yaml.Conditions != nil {
		c = createConditionsFromYaml(yaml.Conditions)
	}
	return NewCondition(yaml.Run, o, c)
}

func createOptionsFromYaml(yamlOptions *map[string]interface{}) *Options {
	options := map[string]interface{}{}

	if yamlOptions != nil {
		options = *yamlOptions
	}
	return NewOptions(options)
}

func createNullableAppSettingsFromYaml(settings *YamlAppSettings) *NullableAppSettings {
	nSettings := NewNullableAppSettings()
	nSettings.AllowFailure = settings.AllowFailure
	nSettings.AnsiColors = settings.AnsiColors
	nSettings.Custom = settings.Custom
	nSettings.FailOnFirstError = settings.FailOnFirstError
	nSettings.GitDirectory = settings.GitDirectory
	nSettings.Includes = settings.Includes
	nSettings.IncludeLevel = settings.IncludeLevel
	nSettings.RunPath = settings.RunPath
	nSettings.RunAsync = settings.RunAsync
	nSettings.Verbosity = settings.Verbosity
	return nSettings
}
