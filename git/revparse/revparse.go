package revparse

import "github.com/captainhook-go/captainhook/git/types"

// AbbrevRef returns a short rev name for the current state
func AbbrevRef(g *types.Cmd) {
	g.AddOption("--abbrev-ref")
}

// ShowTopLevel	return the absolute path to the git repository's root directory
func ShowTopLevel(g *types.Cmd) {
	g.AddOption("--show-toplevel")
}
