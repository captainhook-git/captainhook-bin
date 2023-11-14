package io

import (
	"fmt"
)

type DefaultIO struct {
	verbosity int
	input     Input
}

func NewDefaultIO(verbosity int, arguments map[string]string) *DefaultIO {
	io := DefaultIO{verbosity: verbosity, input: NewStdIn(arguments)}
	return &io
}

func (d *DefaultIO) Verbosity() int {
	return d.verbosity
}

func (d *DefaultIO) Arguments() map[string]string {
	return d.input.Arguments()
}

func (d *DefaultIO) Argument(name, defaultValue string) string {
	return d.input.Argument(name, defaultValue)
}

func (d *DefaultIO) StandardInput() []string {
	return d.input.Data()
}

func (d *DefaultIO) Input() Input {
	return d.input
}

func (d *DefaultIO) IsInteractive() bool {
	return false
}

func (d *DefaultIO) IsQuiet() bool {
	return !(d.verbosity > QUIET)
}

func (d *DefaultIO) IsDebug() bool {
	return d.verbosity == DEBUG
}

func (d *DefaultIO) IsVerbose() bool {
	return d.verbosity == VERBOSE
}

func (d *DefaultIO) Write(message string, newline bool, verbosity int) {
	if d.IsQuiet() {
		return
	}
	var linebreak = ""
	if newline {
		linebreak = "\n"
	}

	if d.verbosity >= verbosity {
		fmt.Print(Colorize(message + linebreak))
	}
}

func (d *DefaultIO) Ask(message string, defaultValue string) string {
	value, err := askForUserInput(message)
	if err != nil {
		value = defaultValue
	}
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}
