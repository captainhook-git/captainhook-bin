package exec

import (
	"encoding/json"
	"errors"
	"github.com/captainhook-go/captainhook/configuration"
	"github.com/captainhook-go/captainhook/git"
	"github.com/captainhook-go/captainhook/info"
	"github.com/captainhook-go/captainhook/io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Initializer creates an empty CaptainHook configuration file
// Is used to set up a empty dummy configuration
type Initializer struct {
	appIO  io.IO
	config string
	force  bool
}

func NewInitializer(appIO io.IO) *Initializer {
	i := Initializer{appIO: appIO, config: "captainhook.json", force: false}
	return &i
}

// UseConfig sets the path to the configuration file
// By default this should be `captainhook.json` but can be changed via the --configuration option.
func (i *Initializer) UseConfig(config string) {
	i.config = config
}

// Force decides if the configuration file will be overwritten without asking
// By default this is false but can be changed with the `--force` option.
func (i *Initializer) Force(force bool) {
	i.force = force
}

// Run executes the Initializer
func (i *Initializer) Run() error {
	i.appIO.Write("Initializing CaptainHook", true, io.VERBOSE)

	gitRoot, gitErr := git.DetectGitDir()
	if gitErr != nil {
		i.appIO.Write("<warning>git repository not found</warning>", true, io.NORMAL)
		return gitErr
	}

	defaultGitDir, detectErr := i.pathToGit(gitRoot)
	if detectErr != nil {
		i.appIO.Write("<warning>could not safely detect git dir please update config manually</warning>", true, io.NORMAL)
	}
	hookConfigs := i.createJsonHookConfigs()

	jsonConfig := &configuration.JsonConfiguration{
		Settings: &configuration.JsonAppSettings{GitDirectory: &defaultGitDir},
		Hooks:    hookConfigs,
	}

	res, jsonErr := json.MarshalIndent(jsonConfig, "", "  ")
	if jsonErr != nil {
		i.appIO.Write("<warning>json generation failed</warning>", true, io.NORMAL)
		return jsonErr
	}

	i.appIO.Write("writing config to '"+i.config+"'", true, io.VERBOSE)
	err := i.writeConfigFile(res)
	if err != nil {
		i.appIO.Write("<warning>initializing failed</warning>", true, io.NORMAL)
		return err
	}
	i.appIO.Write("<ok>successfully initialized</ok>", true, io.NORMAL)
	return nil
}

// pathToGit determines the relative path to the `.git` directory
// If the configuration file is not within the git directory this throws an error
func (i *Initializer) pathToGit(absoluteGit string) (string, error) {
	confDir := path.Dir(i.config)
	absoluteConf, _ := filepath.Abs(confDir)
	if absoluteConf == absoluteGit {
		return ".git", nil
	}

	if !strings.HasPrefix(absoluteConf, absoluteGit) {
		return ".git", errors.New("could not detect .git directory")
	}
	cwdDepth := len(strings.Split(absoluteConf, "/"))
	repoDepth := len(strings.Split(absoluteGit, "/"))

	return "./" + strings.Repeat("../", cwdDepth-repoDepth) + ".git", nil
}

// createJsonHookConfigs create a list of empty hook configurations
// By default this just creates the 3 obvious hooks
func (i *Initializer) createJsonHookConfigs() *map[string]*configuration.JsonHook {
	configs := map[string]*configuration.JsonHook{}

	for _, hook := range []string{info.CommitMsg, info.PreCommit, info.PrePush} {
		configs[hook] = &configuration.JsonHook{
			Actions: make([]*configuration.JsonAction, 0),
		}
	}

	return &configs
}

// writeConfigFile writes the configuration file
// If the file exists it asks the user to confirm to overwrite the existing file
func (i *Initializer) writeConfigFile(res []byte) error {
	doIt := true

	if i.needConfirmationToOverwrite() {
		answer := i.appIO.Ask("<info>"+i.config+"</info> exists! Overwrite? <comment>[y,N]</comment> ", "n")
		doIt = io.AnswerToBool(answer)
	}

	if doIt {
		writeErr := os.WriteFile(i.config, res, 0644)
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}

// needConfirmationToOverwrite returns true is the user has to confirm to overwrite the configuration file
// Only true if the `--force` option is not set and the file exists
func (i *Initializer) needConfirmationToOverwrite() bool {
	return !i.force && io.FileExists(i.config)
}
