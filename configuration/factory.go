package configuration

import (
	"fmt"
	"os"
)

type ConfigFactory interface {
	// CreateConfig creates a CaptainHook configuration
	CreateConfig(path string, cliSettings *NullableAppSettings) (*Configuration, error)
}

func readConfigFile(path string) ([]byte, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("given configuration path is a directory: %s", path)
	}
	jsonData, readErr := os.ReadFile(path)
	if readErr != nil {
		return nil, fmt.Errorf("could not read configuration file at: %s", path)
	}
	return jsonData, nil
}

func copyActionsFromTo(from *Hook, to *Hook) {
	for _, action := range from.GetActions() {
		to.AddAction(action)
	}
}
