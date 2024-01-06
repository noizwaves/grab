package pkg

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const configPath = ".dotlocalbin.yml"

type configRoot struct {
	Binaries []configBinary `yaml:"binaries"`
}

type configBinary struct {
	Name      string                                    `yaml:"name"`
	Source    string                                    `yaml:"source"`
	Platforms map[configPlatformKey]configPlatformValue `yaml:"platforms,omitempty"`
	Version   string                                    `yaml:"version"`
}

type configArchKey = string
type configPlatformKey = string

type configPlatformValue = map[configArchKey]configPlatformArchOverride

type configPlatformArchOverride = [3]string // [platformOveride, archOverride, extOverride]

func parseConfig(path string) (configRoot, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return configRoot{}, fmt.Errorf("Error reading config file: %w", err)
	}

	output := configRoot{}
	err = yaml.Unmarshal(yamlFile, &output)
	if err != nil {
		return configRoot{}, fmt.Errorf("Error parsing config YAML: %w", err)
	}

	return output, nil
}

func loadConfig(homeDir string) (configRoot, error) {
	absPath := path.Join(homeDir, configPath)

	// TODO: validate config here
	return parseConfig(absPath)
}
