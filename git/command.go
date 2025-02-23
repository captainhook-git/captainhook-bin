package git

import (
	"context"
	"github.com/captainhook-go/captainhook/git/types"
	"strings"
)

// command executes a command with a bunch of options
// It uses a Cmd to run the command with the right context and the correct Executor that can be swapped for
// testing purposes.
func command(ctx context.Context, name string, options ...types.Option) (string, error) {
	g := types.NewCmd(name)
	g.AddOptions(options...)
	res, err := g.Exec(ctx, g.Command, g.Debug, g.Options...)

	return strings.TrimSpace(res), err
}

// SetExecutor is the way to substitute the default execution model for testing purposes
func SetExecutor(executor types.Executor) types.Option {
	return func(g *types.Cmd) {
		g.Executor = executor
	}
}

// SetDebug is there to activate debugging for command executions
func SetDebug(debug bool) types.Option {
	return func(g *types.Cmd) {
		g.Debug = debug
	}
}

// Config sets up a `git config` cli command
func Config(options ...types.Option) (string, error) {
	return command(context.Background(), "config", options...)
}

// DiffIndex sets up a `git diff-index` cli command
func DiffIndex(options ...types.Option) (string, error) {
	return command(context.Background(), "diff-index", options...)
}

// DiffTree sets up a `git diff-tree` cli command
func DiffTree(options ...types.Option) (string, error) {
	return command(context.Background(), "diff-tree", options...)
}

// Log sets up a `git log` command
func Log(options ...types.Option) (string, error) {
	return command(context.Background(), "log", options...)
}

// RevParse sets up a `git rev-parse` cli command
func RevParse(options ...types.Option) (string, error) {
	return command(context.Background(), "rev-parse", options...)
}
