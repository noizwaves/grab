package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"runtime"
	"slices"
)

const (
	// defaults are relative to $HOME.

	defaultBinPath       = ".local/bin"
	defaultConfigDirPath = ".grab"

	configFileName    = "config.yml"
	repositoryDirName = "repository"
)

type GrabContext struct {
	Binaries     []*Binary
	BinPath      string
	ConfigPath   string
	Config       *configRoot
	RepoPath     string
	Platform     string
	Architecture string
}

func NewGrabContext(configPathOverride, binPathOverride string) (*GrabContext, error) {
	ctx := context.Background()
	slog.DebugContext(ctx, "Runtime information", "platform", runtime.GOOS, "architecture", runtime.GOARCH)

	configPath, err := getConfigDirPath(configPathOverride)
	if err != nil {
		return nil, fmt.Errorf("error getting config path: %w", err)
	}

	binPath, err := getBinPath(binPathOverride)
	if err != nil {
		return nil, fmt.Errorf("error getting bin path: %w", err)
	}

	configFilePath := path.Join(configPath, configFileName)

	config, err := loadConfig(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	repoPath := path.Join(configPath, repositoryDirName)

	repository, err := loadRepository(repoPath)
	if err != nil {
		return nil, fmt.Errorf("error loading repository: %w", err)
	}

	binaries := make([]*Binary, 0)

	for name, version := range config.Packages {
		located, err := locatePackage(repository, name)
		if err != nil {
			return nil, fmt.Errorf("error locating package information: %w", err)
		}

		binary, err := NewBinary(name, version, *located)
		if err != nil {
			return nil, fmt.Errorf("error constructing binary %q: %w", name, err)
		}

		binaries = append(binaries, binary)
	}

	return &GrabContext{
		Binaries:     binaries,
		BinPath:      binPath,
		ConfigPath:   configFilePath,
		Config:       config,
		RepoPath:     repoPath,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
	}, err
}

func (gc *GrabContext) EnsureBinPathExists() error {
	err := os.MkdirAll(gc.BinPath, 0o755) //nolint:mnd
	if err != nil {
		return fmt.Errorf("error creating bin path directory: %w", err)
	}

	return nil
}

func (gc *GrabContext) SavePackage(packageConfig *ConfigPackage) (string, error) {
	packagePath := path.Join(gc.RepoPath, packageConfig.Metadata.Name+".yml")

	err := savePackage(packageConfig, packagePath)
	if err != nil {
		return "", fmt.Errorf("error saving package: %w", err)
	}

	return packagePath, nil
}

func getPackageNames(repository *repository) []string {
	names := make([]string, len(repository.Packages))
	for idx, pkg := range repository.Packages {
		names[idx] = pkg.Metadata.Name
	}

	return names
}

func locatePackage(repository *repository, name string) (*ConfigPackage, error) {
	ctx := context.Background()
	slog.DebugContext(ctx, "Looking for configured package in repository", "name", name)

	idx := slices.IndexFunc(repository.Packages, func(p *ConfigPackage) bool { return p.Metadata.Name == name })
	if idx == -1 {
		ctx := context.Background()
		slog.ErrorContext(ctx, "Package missing from repository", "name", name)

		slog.DebugContext(ctx, "Repository contains", "packageNames", getPackageNames(repository))

		return nil, fmt.Errorf("package %q missing from repository", name)
	}

	return repository.Packages[idx], nil
}

func getConfigDirPath(override string) (string, error) {
	var configDirPath string
	if override != "" {
		configDirPath = override
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error determining home directory: %w", err)
		}

		configDirPath = path.Join(homeDir, defaultConfigDirPath)
	}

	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("config directory does not exist: %w", err)
	}

	return configDirPath, nil
}

func getBinPath(override string) (string, error) {
	if override != "" {
		return override, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error determining home directory: %w", err)
	}

	binPath := path.Join(homeDir, defaultBinPath)

	return binPath, nil
}
