package commands

import (
	"errors"
	"fmt"
	"github.com/captainhook-go/captainhook/configuration"
	"github.com/captainhook-go/captainhook/info"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

// configurationAware is for all commands that need to read or write a configuration
func configurationAware(cmd *cobra.Command) {
	var configPath = info.Config
	cmd.Flags().StringP("configuration", "c", configPath, "path to your CaptainHook config")
}

// repositoryAware is for all commands that need a repository to execute git commands or access the .git directory
func repositoryAware(cmd *cobra.Command) {
	var repoPath = ""
	cmd.Flags().StringP("git-directory", "g", repoPath, "path to your .git directory")
}

// setUpConfig creates a Configuration struct
// It uses a JsonAppSettings struct because it has nullable properties.
// This way we can figure out what settings were actually set and which were not later.
// This is important since the command line options should supersede all other ways of
// configuring the Cap'n.
func setUpConfig(cmd *cobra.Command, fileRequired bool) (*configuration.Configuration, error) {
	nullableSettings := &configuration.NullableAppSettings{}

	detectColor(cmd, nullableSettings)
	detectGitDir(cmd, nullableSettings)
	detectVerbosity(cmd, nullableSettings)

	confPath := info.Config
	confOption, _ := cmd.Flags().GetString("configuration")
	if len(confOption) > 0 {
		confPath = confOption
	}
	var factory configuration.ConfigFactory
	ext := strings.ToLower(filepath.Ext(confPath))
	switch ext {
	case ".yaml", ".yml":
		factory = configuration.NewYamlFactory()
	case ".json":
		factory = configuration.NewJsonFactory()
	default:
		return nil, errors.New(fmt.Sprintf("Unknown configuration format: %s", ext))
	}

	conf, confErr := factory.CreateConfig(confPath, nullableSettings)
	if confErr != nil {
		return nil, confErr
	}
	if fileRequired && !conf.IsLoadedFromFile() {
		return nil, errors.New("configuration file not found")
	}
	return conf, nil
}

// detectColor is checking the `--no-color` option and sets the configuration accordingly
func detectColor(cmd *cobra.Command, settings *configuration.NullableAppSettings) {
	noColor := getNoColor(cmd)
	if noColor {
		falsePointer := false
		settings.AnsiColors = &falsePointer
	}
}

func getNoColor(cmd *cobra.Command) bool {
	noColor, _ := cmd.Flags().GetBool("no-color")
	return noColor
}

// detectGitDir is checking the `--git-directory` option and sets the configuration accordingly
func detectGitDir(cmd *cobra.Command, settings *configuration.NullableAppSettings) {
	repoOption := getGitDir(cmd)
	if len(repoOption) > 0 {
		settings.GitDirectory = &repoOption
	}
}

func getGitDir(cmd *cobra.Command) string {
	repoPath := ""
	repoOption, _ := cmd.Flags().GetString("git-directory")
	if len(repoOption) > 0 {
		repoPath = repoOption
	}
	return repoPath
}

// detectVerbosity is checking the `--verbose` and `--debug` options and sets the configuration accordingly
func detectVerbosity(cmd *cobra.Command, settings *configuration.NullableAppSettings) {
	v := getVerbosity(cmd)
	if v != "normal" {
		settings.Verbosity = &v
	}
}

func getVerbosity(cmd *cobra.Command) string {
	verbosity := "normal"
	quiet, _ := cmd.Flags().GetBool("quiet")
	if quiet {
		verbosity = "quiet"
	}
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		verbosity = "verbose"
	}
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		verbosity = "debug"
	}
	return verbosity
}

// mapArgs is mapping the command arguments into a named map
// To not require the knowledge on which position an argument came in everywhere arguments can now be accessed by name.
// Instead of `.getArgument(0)` wen can write `.getArgument("message-file")`
// In addition the current command gets added again under the index "command".
func mapArgs(names []string, args []string, cmd string) map[string]string {
	m := make(map[string]string)
	for index, name := range names {
		if len(args) > index {
			m[name] = args[index]
		}
	}
	m[info.ArgCommand] = cmd
	return m
}
