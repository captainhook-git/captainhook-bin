package configuration

import (
	"github.com/captainhook-go/captainhook/info"
	"github.com/captainhook-go/captainhook/io"
)

type Configuration struct {
	size       int64
	path       string
	fileExists bool
	settings   *AppSettings
	hooks      map[string]*Hook
}

func NewConfiguration(path string, fileExists bool) *Configuration {
	c := &Configuration{path: path, fileExists: fileExists, settings: NewDefaultAppSettings()}
	c.init()
	return c
}

func (c *Configuration) init() {
	c.hooks = map[string]*Hook{}
	for _, hook := range info.GetValidHooks() {
		c.hooks[hook] = NewHook(hook, false)
	}
}

func (c *Configuration) IsLoadedFromFile() bool {
	return c.fileExists
}

func (c *Configuration) IsHookEnabled(hook string) bool {
	return c.HookConfig(hook).IsEnabled()
}

func (c *Configuration) Path() string {
	return c.path
}

func (c *Configuration) RunPath() string {
	return c.settings.RunPath
}

func (c *Configuration) CustomSettings() map[string]string {
	return c.settings.Custom
}

func (c *Configuration) GitDirectory() string {
	gitDir := c.settings.GitDirectory
	if len(gitDir) < 1 {
		gitDir = ".git"
	}
	return gitDir
}

func (c *Configuration) AnsiColors() bool {
	return c.settings.AnsiColors
}

func (c *Configuration) Verbosity() int {
	return MapVerbosity(c.settings.Verbosity)
}

func (c *Configuration) IsFailureAllowed() bool {
	return c.settings.AllowFailure
}

func (c *Configuration) FailOnFirstError() bool {
	return c.settings.FailOnFirstError
}

func (c *Configuration) RunAsync() bool {
	return c.settings.RunAsync
}

func (c *Configuration) Includes() []string {
	return c.settings.Includes
}

func (c *Configuration) MaxIncludeLevel() int {
	return c.settings.IncludeLevel
}

func (c *Configuration) HookConfig(hook string) *Hook {
	return c.hooks[hook]
}

// overwriteSettings will overwrite every setting that is set in the jsonConfig.
func (c *Configuration) overwriteSettings(settings *NullableAppSettings) {
	if settings == nil {
		return
	}

	if settings.AllowFailure != nil {
		c.settings.AllowFailure = *settings.AllowFailure
	}
	if settings.AnsiColors != nil {
		c.settings.AnsiColors = *settings.AnsiColors
	}
	if (settings.Custom) != nil {
		c.settings.Custom = *settings.Custom
	}
	if settings.FailOnFirstError != nil {
		c.settings.FailOnFirstError = *settings.FailOnFirstError
	}
	if settings.GitDirectory != nil {
		c.settings.GitDirectory = *settings.GitDirectory
	}
	if settings.RunPath != nil {
		c.settings.RunPath = *settings.RunPath
	}
	if settings.RunAsync != nil {
		c.settings.RunAsync = *settings.RunAsync
	}
	if settings.Verbosity != nil {
		c.settings.Verbosity = *settings.Verbosity
	}
}

func MapVerbosity(verbosity string) int {
	verbosityMap := map[string]int{
		"quiet":   io.QUIET,
		"normal":  io.NORMAL,
		"verbose": io.VERBOSE,
		"debug":   io.DEBUG,
	}
	verbosityIO, ok := verbosityMap[verbosity]
	if !ok {
		verbosityIO = io.NORMAL
	}
	return verbosityIO
}

func UnMapVerbosity(verbosity int) string {
	verbosityMap := map[int]string{
		io.QUIET:   "quiet",
		io.NORMAL:  "normal",
		io.VERBOSE: "verbose",
		io.DEBUG:   "debug",
	}
	verbosityConfig, ok := verbosityMap[verbosity]
	if !ok {
		verbosityConfig = "normal"
	}
	return verbosityConfig
}
