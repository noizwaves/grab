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
	Name    string `yaml:"name"`
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
}

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

func loadConfig() (configRoot, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return configRoot{}, fmt.Errorf("Error determining home directory: %w", err)
	}

	absPath := path.Join(homeDir, configPath)

	return parseConfig(absPath)
}
