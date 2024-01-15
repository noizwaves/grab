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
	Packages []configPackage
}

type configPackage struct {
	APIVersion string                `yaml:"apiVersion"`
	Kind       string                `yaml:"kind"`
	Metadata   configPackageMetadata `yaml:"metadata"`
	Spec       configPackageSpec     `yaml:"spec"`
}

type configPackageMetadata struct {
	Name string `yaml:"name"`
}

type configPackageSpec struct {
	GitHubRelease configGitHubRelease `yaml:"gitHubRelease"`
	Program       configProgram       `yaml:"program"`
}

type configGitHubRelease struct {
	Org          string            `yaml:"org"`
	Repo         string            `yaml:"repo"`
	Name         string            `yaml:"name"`
	VersionRegex string            `yaml:"versionRegex"`
	FileName     map[string]string `yaml:"fileName"`
}

type configProgram struct {
	VersionArgs  []string `yaml:"versionArgs"`
	VersionRegex string   `yaml:"versionRegex"`
}

func loadConfig(path string) (configRoot, error) {
	slog.Info("Loading config file from disk", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return configRoot{}, fmt.Errorf("error reading config file: %w", err)
	}

	slog.Debug("Loaded config from disk", "content", string(data))

	output := configRoot{}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	err = decoder.Decode(&output)

	if err != nil {
		return configRoot{}, fmt.Errorf("error parsing config YAML: %w", err)
	}

	return output, nil
}

func saveConfig(config *configRoot, path string) error {
	slog.Info("Saving config file to disk", "path", path)
	var buf bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buf)
	yamlEncoder.SetIndent(2) //nolint:gomnd
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

func loadPackage(path string) (configPackage, error) {
	slog.Info("Loading package config from disk", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return configPackage{}, fmt.Errorf("error reading package file: %w", err)
	}

	slog.Debug("Loaded package config from disk", "content", string(data))

	output := configPackage{}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	err = decoder.Decode(&output)
	if err != nil {
		return configPackage{}, fmt.Errorf("error parsing package YAML: %w", err)
	}

	return output, nil
}

func loadRepository(repoPath string) (repository, error) {
	slog.Info("Loading packages from repository on disk", "repoPath", repoPath)
	packages := []configPackage{}

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
		return repository{}, fmt.Errorf("error loading repository: %w", err)
	}

	slog.Info("Loaded packages from repository", "count", len(packages))

	return repository{
		Packages: packages,
	}, nil
}
