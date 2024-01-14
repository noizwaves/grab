package pkg

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"runtime"
	"slices"
)

const (
	localBinPath   = ".local/bin"
	configPath     = ".garb/config.yml"
	repositoryPath = ".garb/repository"
)

type Context struct {
	Binaries     []Binary
	BinPath      string
	ConfigPath   string
	Config       *configRoot
	Platform     string
	Architecture string
}

func getPackageNames(repository *repository) []string {
	names := make([]string, len(repository.Packages))
	for idx, pkg := range repository.Packages {
		names[idx] = pkg.Metadata.Name
	}

	return names
}

func locatePackage(repository *repository, name string) (*configPackage, error) {
	slog.Debug("Looking for configured package in repository", "name", name)
	idx := slices.IndexFunc(repository.Packages, func(p configPackage) bool { return p.Metadata.Name == name })
	if idx == -1 {
		slog.Error("Package missing from repository", "name", name)

		slog.Debug("Repository contains", "packageNames", getPackageNames(repository))

		return nil, fmt.Errorf("package %q missing from repository", name)
	}

	return &repository.Packages[idx], nil
}

func NewContext() (Context, error) {
	slog.Debug("Runtime information", "platform", runtime.GOOS, "architecture", runtime.GOARCH)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Context{}, fmt.Errorf("error determining home directory: %w", err)
	}
	slog.Debug("Configuration information", "homeDir", homeDir)

	configPath := path.Join(homeDir, configPath)
	config, err := loadConfig(configPath)
	if err != nil {
		return Context{}, fmt.Errorf("error loading config: %w", err)
	}

	repoPath := path.Join(homeDir, repositoryPath)
	repository, err := loadRepository(repoPath)
	if err != nil {
		return Context{}, fmt.Errorf("error loading repository: %w", err)
	}

	binaries := make([]Binary, 0)
	for name, version := range config.Packages {
		located, err := locatePackage(&repository, name)
		if err != nil {
			return Context{}, fmt.Errorf("error locating package information: %w", err)
		}

		binary, err := NewBinary(name, version, *located)
		if err != nil {
			return Context{}, fmt.Errorf("error constructing binary %q: %w", name, err)
		}

		binaries = append(binaries, binary)
	}

	return Context{
		Binaries:     binaries,
		BinPath:      path.Join(homeDir, localBinPath),
		ConfigPath:   configPath,
		Config:       &config,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
	}, err
}
