package configuration

import (
	"errors"
	"fmt"
	"github.com/captainhook-go/captainhook/info"
	"github.com/captainhook-go/captainhook/io"
	"gopkg.in/yaml.v3"
	"log"
	"path/filepath"
)

type YamlFactory struct {
	includeLevel    int
	maxIncludeLevel int
}

// CreateConfig creates a default configuration in case the file exists it is loaded
func (f *YamlFactory) CreateConfig(path string, cliSettings *NullableAppSettings) (*Configuration, error) {
	c := NewConfiguration(path, io.FileExists(path))
	if c.fileExists {
		cErr := f.loadFromFile(c)
		if cErr != nil {
			return c, cErr
		}
	}
	// load the local config "captainhook.config.yml"
	sErr := f.loadSettingsFile(c)
	if sErr != nil {
		return c, sErr
	}
	// everything provided from the command line should overwrite any loaded configuration
	// this works even if there is an error because then you have a default configuration
	c.overwriteSettings(cliSettings)
	return c, nil
}

// setupConfig creates a new configuration and loads the json file if it exists
func (f *YamlFactory) setupConfig(path string) (*Configuration, error) {
	var err error
	c := NewConfiguration(path, io.FileExists(path))
	if c.fileExists {
		err = f.loadFromFile(c)
	}
	return c, err
}

func (f *YamlFactory) loadFromFile(c *Configuration) error {
	yamlBytes, readError := readConfigFile(c.path)
	if readError != nil {
		return readError
	}
	configurationYaml, decodeErr := f.decodeConfigYaml(yamlBytes)
	if decodeErr != nil {
		return fmt.Errorf("unable to parse yaml: %s %s", c.path, decodeErr.Error())
	}
	c.overwriteSettings(createNullableAppSettingsFromYaml(configurationYaml.Settings))

	if configurationYaml.Hooks == nil {
		return errors.New("no hooks config found")
	}
	includeErr := f.appendIncludedConfiguration(c)
	if includeErr != nil {
		return includeErr
	}

	for hookName, hookConfigYaml := range *configurationYaml.Hooks {
		hookConfig := c.HookConfig(hookName)
		hookConfig.isEnabled = true
		for _, actionYaml := range hookConfigYaml.Actions {
			if !f.isValidAction(actionYaml) {
				return fmt.Errorf("invalid action config in %s", hookName)
			}
			hookConfig.AddAction(createActionFromYaml(actionYaml))
		}
	}
	return nil
}

func (f *YamlFactory) loadSettingsFile(c *Configuration) error {
	directory := filepath.Dir(c.path)
	filePath := directory + "/captainhook.config.yaml"

	// no local config file to load just exit
	if !io.FileExists(filePath) {
		return nil
	}

	yamlBytes, readError := readConfigFile(filePath)
	if readError != nil {
		return readError
	}
	appSettingsYaml, decodeErr := f.decodeSettingYaml(yamlBytes)
	if decodeErr != nil {
		return fmt.Errorf("unable to parse json: %s %s", filePath, decodeErr.Error())
	}
	// overwrite current settings
	c.overwriteSettings(createNullableAppSettingsFromYaml(appSettingsYaml))
	return nil
}

func (f *YamlFactory) appendIncludedConfiguration(c *Configuration) error {
	f.detectMaxIncludeLevel(c)
	if f.includeLevel < f.maxIncludeLevel {
		f.includeLevel++
		includes, err := f.loadIncludedConfigs(c.Includes(), c.path)
		if err != nil {
			return err
		}
		for _, configToInclude := range includes {
			f.mergeHookConfigs(configToInclude, c)
		}
		f.includeLevel--
	}
	return nil
}

func (f *YamlFactory) mergeHookConfigs(from, to *Configuration) {
	for _, hook := range info.GetValidHooks() {
		// This `Enable` is solely to overwrite the main configuration in the special case that the hook
		// is not configured at all. In this case the empty config is disabled by default, and adding an
		// empty hook config just to enable the included actions feels a bit dull.
		// Since the main hook is processed last (if one is configured) the enabled flag will be overwritten
		// once again by the main config value. This is to make sure that if somebody disables a hook in its
		// main configuration no actions will get executed, even if we have enabled hooks in any include file.
		targetHookConfig := to.HookConfig(hook)
		targetHookConfig.Enable()
		copyActionsFromTo(from.HookConfig(hook), targetHookConfig)
	}
}

func (f *YamlFactory) decodeConfigYaml(yamlInBytes []byte) (YamlConfiguration, error) {
	var yConfig YamlConfiguration
	if err := yaml.Unmarshal(yamlInBytes, &yConfig); err != nil {
		log.Fatalf("could not load yaml to struct: %v", err)
	}
	return yConfig, nil
}

func (f *YamlFactory) decodeSettingYaml(yamlAsBytes []byte) (*YamlAppSettings, error) {
	var settings YamlAppSettings
	marshalError := yaml.Unmarshal(yamlAsBytes, &settings)
	if marshalError != nil {
		return nil, fmt.Errorf("could not load yaml to struct: %s", marshalError.Error())
	}
	return &settings, nil
}

func (f *YamlFactory) detectMaxIncludeLevel(c *Configuration) {
	// read the include-level setting only for the actual configuration not any included ones
	if f.includeLevel == 0 {
		f.maxIncludeLevel = c.MaxIncludeLevel()
	}
}

func (f *YamlFactory) loadIncludedConfigs(includes []string, path string) ([]*Configuration, error) {
	var configs []*Configuration
	directory := filepath.Dir(path)

	for _, file := range includes {
		config, err := f.includeConfig(directory + "/" + file)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}

func (f *YamlFactory) includeConfig(path string) (*Configuration, error) {
	if !io.FileExists(path) {
		return nil, fmt.Errorf("config to include not found: %s", path)
	}
	return f.setupConfig(path)
}

func (f *YamlFactory) isValidAction(action YamlAction) bool {
	if len(action.Run) == 0 {
		return false
	}
	for _, condition := range action.Conditions {
		if len(condition.Run) == 0 {
			return false
		}
	}
	return true
}

func NewYamlFactory() *YamlFactory {
	return &YamlFactory{includeLevel: 0}
}
