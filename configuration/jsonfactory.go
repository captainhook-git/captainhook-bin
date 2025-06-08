package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/captainhook-go/captainhook/info"
	"github.com/captainhook-go/captainhook/io"
	"path/filepath"
)

type JsonFactory struct {
	includeLevel    int
	maxIncludeLevel int
}

// CreateConfig creates a default configuration in case the file exists it is loaded
func (f *JsonFactory) CreateConfig(path string, cliSettings *NullableAppSettings) (*Configuration, error) {
	c, cErr := f.setupConfig(path)
	if cErr != nil {
		return c, cErr
	}
	// load the local config "captainhook.config.json"
	sErr := f.loadSettingsFile(c)
	if sErr != nil {
		return c, sErr
	}
	// everything provided from the command line should overwrite any loaded configuration
	c.overwriteSettings(cliSettings)
	return c, nil
}

// setupConfig creates a new configuration and loads the json file if it exists
func (f *JsonFactory) setupConfig(path string) (*Configuration, error) {
	var err error
	c := NewConfiguration(path, io.FileExists(path))
	if c.fileExists {
		err = f.loadFromFile(c)
	}
	return c, err
}

func (f *JsonFactory) loadFromFile(c *Configuration) error {
	jsonBytes, readError := readConfigFile(c.path)
	if readError != nil {
		return readError
	}
	configurationJson, decodeErr := f.decodeConfigJson(jsonBytes)
	if decodeErr != nil {
		return fmt.Errorf("unable to parse json: %s %s", c.path, decodeErr.Error())
	}
	c.overwriteSettings(createNullableAppSettingsFromJson(configurationJson.Settings))

	if configurationJson.Hooks == nil {
		return errors.New("no hooks config found")
	}
	includeErr := f.appendIncludedConfiguration(c)
	if includeErr != nil {
		return includeErr
	}

	for hookName, hookConfigJson := range *configurationJson.Hooks {
		hookConfig := c.HookConfig(hookName)
		hookConfig.isEnabled = true
		for _, actionJson := range hookConfigJson.Actions {
			if !f.isValidAction(actionJson) {
				return fmt.Errorf("invalid action config in %s", hookName)
			}
			hookConfig.AddAction(createActionFromJson(actionJson))
		}
	}
	return nil
}

func (f *JsonFactory) loadSettingsFile(c *Configuration) error {
	directory := filepath.Dir(c.path)
	filePath := directory + "/captainhook.config.json"

	// no local config file to load just exit
	if !io.FileExists(filePath) {
		return nil
	}

	jsonBytes, readError := readConfigFile(filePath)
	if readError != nil {
		return readError
	}
	appSettingJson, decodeErr := f.decodeSettingJson(jsonBytes)
	if decodeErr != nil {
		return fmt.Errorf("unable to parse json: %s %s", filePath, decodeErr.Error())
	}
	// overwrite current settings
	c.overwriteSettings(createNullableAppSettingsFromJson(appSettingJson))
	return nil
}

func (f *JsonFactory) appendIncludedConfiguration(c *Configuration) error {
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

func (f *JsonFactory) mergeHookConfigs(from, to *Configuration) {
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

func (f *JsonFactory) decodeConfigJson(jsonInBytes []byte) (*JsonConfiguration, error) {
	var jConfig *JsonConfiguration
	if !json.Valid(jsonInBytes) {
		return nil, fmt.Errorf("json configuration is invalid")
	}
	marshalError := json.Unmarshal(jsonInBytes, &jConfig)
	if marshalError != nil {
		return nil, fmt.Errorf("could not load json to struct: %s", marshalError.Error())
	}
	return jConfig, nil
}

func (f *JsonFactory) decodeSettingJson(jsonInBytes []byte) (*JsonAppSettings, error) {
	var jSettings JsonAppSettings
	if !json.Valid(jsonInBytes) {
		return nil, fmt.Errorf("json configuration is invalid")
	}
	marshalError := json.Unmarshal(jsonInBytes, &jSettings)
	if marshalError != nil {
		return nil, fmt.Errorf("could not load json to struct: %s", marshalError.Error())
	}
	return &jSettings, nil
}

func (f *JsonFactory) detectMaxIncludeLevel(c *Configuration) {
	// read the include-level setting only for the actual configuration not any included ones
	if f.includeLevel == 0 {
		f.maxIncludeLevel = c.MaxIncludeLevel()
	}
}

func (f *JsonFactory) loadIncludedConfigs(includes []string, path string) ([]*Configuration, error) {
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

func (f *JsonFactory) includeConfig(path string) (*Configuration, error) {
	if !io.FileExists(path) {
		return nil, fmt.Errorf("config to include not found: %s", path)
	}
	return f.setupConfig(path)
}

func (f *JsonFactory) isValidAction(actionJson *JsonAction) bool {
	if len(actionJson.Run) == 0 {
		return false
	}
	for _, condition := range actionJson.Conditions {
		if len(condition.Run) == 0 {
			return false
		}
	}
	return true
}

func NewJsonFactory() *JsonFactory {
	return &JsonFactory{includeLevel: 0}
}
