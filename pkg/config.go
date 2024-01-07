package pkg

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type configRoot struct {
	Binaries []configBinary `yaml:"binaries"`
}

type configBinary struct {
	Name         string                                    `yaml:"name"`
	Source       string                                    `yaml:"source"`
	Platforms    map[configPlatformKey]configPlatformValue `yaml:"platforms,omitempty"`
	Version      string                                    `yaml:"version"`
	VersionFlags []string                                  `yaml:"versionFlags"`
	VersionRegex string                                    `yaml:"versionRegex"`
}

type (
	configArchKey     = string
	configPlatformKey = string
)

type configPlatformValue = map[configArchKey]configPlatformArchOverride

type configPlatformArchOverride = [3]string // [platformOveride, archOverride, extOverride]

func parseConfig(path string) (configRoot, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return configRoot{}, fmt.Errorf("error reading config file: %w", err)
	}

	output := configRoot{}
	err = yaml.Unmarshal(yamlFile, &output)
	if err != nil {
		return configRoot{}, fmt.Errorf("error parsing config YAML: %w", err)
	}

	return output, nil
}

func loadConfig(homeDir string) (configRoot, error) {
	absPath := path.Join(homeDir, configPath)

	return parseConfig(absPath)
}
