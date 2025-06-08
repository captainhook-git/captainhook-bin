package configuration

type Action struct {
	run        string
	settings   *ActionSettings
	conditions []*Condition
	options    *Options
}

func (a *Action) Run() string {
	return a.run
}

func (a *Action) Options() *Options {
	return a.options
}

func (a *Action) Label() string {
	if len(a.settings.Label) > 0 {
		return a.settings.Label
	}
	return a.Run()
}

func (a *Action) IsFailureAllowed() bool {
	return a.settings.AllowFailure
}

func (a *Action) WorkingDir() string {
	return a.settings.WorkingDir
}

func (a *Action) Conditions() []*Condition {
	return a.conditions
}
