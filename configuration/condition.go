package configuration

type Condition struct {
	run        string
	options    *Options
	conditions []*Condition
}

func (c *Condition) Run() string {
	return c.run
}

func (c *Condition) Options() *Options {
	return c.options
}

func (c *Condition) Conditions() []*Condition {
	return c.conditions
}

func NewCondition(cmd string, o *Options, c []*Condition) *Condition {
	return &Condition{run: cmd, options: o, conditions: c}
}
