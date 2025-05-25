package config

import "github.com/captainhook-go/captainhook/git/types"

// Get returns the value of the requested config value
func Get(name string) func(*types.Cmd) {
	return func(g *types.Cmd) {
		g.AddOption("--get")
		g.AddOption(name)
	}
}
