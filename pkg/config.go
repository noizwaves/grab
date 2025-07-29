package pkg

import (
	"bytes"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type configRoot struct {
	Packages map[string]string `yaml:"packages"`
}

type repository struct {
	Packages []*ConfigPackage
}

type ConfigPackage struct {
	APIVersion string                `yaml:"apiVersion"`
	Kind       string                `yaml:"kind"`
	Metadata   ConfigPackageMetadata `yaml:"metadata"`
	Spec       ConfigPackageSpec     `yaml:"spec"`
}

type ConfigPackageMetadata struct {
	Name string `yaml:"name"`
}

type ConfigPackageSpec struct {
	GitHubRelease ConfigGitHubRelease `yaml:"gitHubRelease"`
	Program       ConfigProgram       `yaml:"program"`
}

type ConfigGitHubRelease struct {
	Org                string            `yaml:"org"`
	Repo               string            `yaml:"repo"`
	Name               string            `yaml:"name"`
	VersionRegex       string            `yaml:"versionRegex"`
	FileName           map[string]string `yaml:"fileName"`
	EmbeddedBinaryPath map[string]string `yaml:"embeddedBinaryPath,omitempty"`
}

type ConfigProgram struct {
	VersionArgs  []string `yaml:"versionArgs,flow"`
	VersionRegex string   `yaml:"versionRegex"`
}

func loadConfig(path string) (*configRoot, error) {
	slog.Info("Loading config file from disk", "path", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	slog.Debug("Loaded config from disk", "content", string(data))

	output := configRoot{}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	err = decoder.Decode(&output)
	if err != nil {
		return nil, fmt.Errorf("error parsing config YAML: %w", err)
	}

	return &output, nil
}

func saveConfig(config *configRoot, path string) error {
	slog.Info("Saving config file to disk", "path", path)

	var buf bytes.Buffer

	yamlEncoder := yaml.NewEncoder(&buf)
	yamlEncoder.SetIndent(2) //nolint:mnd

	err := yamlEncoder.Encode(config)
	if err != nil {
		return fmt.Errorf("error serializing config: %w", err)
	}

	data := buf.Bytes()

	slog.Debug("Writing config to disk", "content", string(data))

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error opening existing config file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}

func loadPackage(path string) (*ConfigPackage, error) {
	slog.Info("Loading package config from disk", "path", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading package file: %w", err)
	}

	slog.Debug("Loaded package config from disk", "content", string(data))

	output := ConfigPackage{}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	err = decoder.Decode(&output)
	if err != nil {
		return nil, fmt.Errorf("error parsing package YAML: %w", err)
	}

	return &output, nil
}

func savePackage(packageConfig *ConfigPackage, path string) error {
	slog.Info("Saving package config to disk", "path", path)

	var buf bytes.Buffer

	yamlEncoder := yaml.NewEncoder(&buf)
	yamlEncoder.SetIndent(2) //nolint:mnd

	err := yamlEncoder.Encode(packageConfig)
	if err != nil {
		return fmt.Errorf("error serializing package config: %w", err)
	}

	data := buf.Bytes()

	slog.Debug("Writing package config to disk", "content", string(data))

	err = os.WriteFile(path, data, 0o644) //nolint:gosec,mnd
	if err != nil {
		return fmt.Errorf("error writing package config: %w", err)
	}

	return nil
}

func loadRepository(repoPath string) (*repository, error) {
	slog.Info("Loading packages from repository on disk", "repoPath", repoPath)

	packages := []*ConfigPackage{}

	err := filepath.Walk(repoPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error reading file system: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".yml") {
			slog.Debug("Skipping non .yml file", "path", path)

			return nil
		}

		loaded, err := loadPackage(path)
		if err != nil {
			return fmt.Errorf("error loading package config: %w", err)
		}

		packages = append(packages, loaded)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error loading repository: %w", err)
	}

	slog.Info("Loaded packages from repository", "count", len(packages))

	return &repository{
		Packages: packages,
	}, nil
}
