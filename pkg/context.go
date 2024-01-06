package pkg

import (
	"fmt"
	"os"
	"runtime"
)

type Context struct {
	Binaries     []Binary
	HomeDir      string
	Platform     string
	Architecture string
}

func NewContext() (Context, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Context{}, fmt.Errorf("Error determining home directory: %w", err)
	}

	config, err := loadConfig(homeDir)
	if err != nil {
		return Context{}, fmt.Errorf("Error loading config: %w", err)
	}

	binaries := make([]Binary, 0)
	for _, cb := range config.Binaries {
		binary, err := NewBinary(cb)
		if err != nil {
			return Context{}, fmt.Errorf("Error constructing binary %q: %w", cb.Name, err)
		}

		binaries = append(binaries, binary)
	}

	return Context{
		Binaries:     binaries,
		HomeDir:      homeDir,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
	}, err
}
