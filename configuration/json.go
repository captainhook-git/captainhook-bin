package configuration

type JsonConfiguration struct {
	Settings *JsonAppSettings      `json:"config,omitempty"`
	Hooks    *map[string]*JsonHook `json:"hooks,omitempty"`
}

type JsonHook struct {
	Actions []*JsonAction `json:"actions"`
}

type JsonAction struct {
	Run        string                  `json:"run,omitempty"`
	Conditions []*JsonCondition        `json:"conditions,omitempty"`
	Options    *map[string]interface{} `json:"options,omitempty"`
	Settings   *JsonActionSettings     `json:"config,omitempty"`
}

type JsonActionSettings struct {
	Label        *string `json:"label,omitempty"`
	AllowFailure *bool   `json:"allow-failure,omitempty"`
	RunAsync     *bool   `json:"run-async,omitempty"`
	WorkingDir   *string `json:"working-dir,omitempty"`
}

type JsonCondition struct {
	Run        string                  `json:"run,omitempty"`
	Options    *map[string]interface{} `json:"options,omitempty"`
	Conditions []*JsonCondition        `json:"conditions,omitempty"`
}

type JsonAppSettings struct {
	AllowFailure     *bool              `json:"allow-failure,omitempty"`
	AnsiColors       *bool              `json:"ansi-colors,omitempty"`
	Custom           *map[string]string `json:"custom,omitempty"`
	FailOnFirstError *bool              `json:"fail-on-first-error,omitempty"`
	GitDirectory     *string            `json:"git-directory,omitempty"`
	Includes         *[]string          `json:"includes,omitempty"`
	IncludeLevel     *int               `json:"includes-level,omitempty"`
	RunPath          *string            `json:"run-path,omitempty"`
	RunAsync         *bool              `json:"run-async,omitempty"`
	Verbosity        *string            `json:"verbosity,omitempty"`
}

func createActionFromJson(json *JsonAction) *Action {
	return &Action{
		run:        json.Run,
		settings:   createActionSettingsFromJson(json.Settings),
		conditions: createConditionsFromJson(json.Conditions),
		options:    createOptionsFromJson(json.Options),
	}
}

func createActionSettingsFromJson(json *JsonActionSettings) *ActionSettings {
	a := NewDefaultActionSettings()
	if json == nil {
		return a
	}
	if json.AllowFailure != nil {
		a.AllowFailure = *json.AllowFailure
	}
	if json.WorkingDir != nil {
		a.WorkingDir = *json.WorkingDir
	}
	if json.Label != nil {
		a.Label = *json.Label
	}
	return a
}

func createConditionsFromJson(jsonConditions []*JsonCondition) []*Condition {
	var conditions []*Condition
	if jsonConditions == nil {
		return conditions
	}
	for _, condition := range jsonConditions {
		conditions = append(conditions, createConditionFromJson(condition))
	}
	return conditions
}

func createConditionFromJson(json *JsonCondition) *Condition {
	var c []*Condition

	// default value empty options
	opts := map[string]interface{}{}
	o := NewOptions(opts)

	if json.Options != nil {
		o = createOptionsFromJson(json.Options)
	}
	if json.Conditions != nil {
		c = createConditionsFromJson(json.Conditions)
	}
	return NewCondition(json.Run, o, c)
}

func createOptionsFromJson(jsonOptions *map[string]interface{}) *Options {
	options := map[string]interface{}{}

	if jsonOptions != nil {
		options = *jsonOptions
	}
	return NewOptions(options)
}

func createNullableAppSettingsFromJson(appSettingJson *JsonAppSettings) *NullableAppSettings {
	nSettings := NewNullableAppSettings()
	nSettings.AllowFailure = appSettingJson.AllowFailure
	nSettings.AnsiColors = appSettingJson.AnsiColors
	nSettings.Custom = appSettingJson.Custom
	nSettings.FailOnFirstError = appSettingJson.FailOnFirstError
	nSettings.GitDirectory = appSettingJson.GitDirectory
	nSettings.Includes = appSettingJson.Includes
	nSettings.IncludeLevel = appSettingJson.IncludeLevel
	nSettings.RunPath = appSettingJson.RunPath
	nSettings.RunAsync = appSettingJson.RunAsync
	nSettings.Verbosity = appSettingJson.Verbosity
	return nSettings
}
