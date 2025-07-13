package pkg

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeEmptyConfig(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	err := os.MkdirAll(path.Join(dir, ".grab", "repository"), 0o755)
	if err != nil {
		t.Fatal(err)
	}

	emptyConfig := []byte(`packages: {}`)

	err = os.WriteFile(path.Join(dir, ".grab", "config.yml"), emptyConfig, 0o644) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

func TestNewContext(t *testing.T) {
	// Use simple as $HOME
	t.Setenv("HOME", "../testdata/simple")

	// Set up an overrides directory
	emptyConfigDir := makeEmptyConfig(t)

	t.Run("NoOverrides", func(t *testing.T) {
		result, err := NewContext("", "")

		assert.NoError(t, err)

		assert.Equal(t, result.ConfigPath, path.Join("../testdata/simple/.grab/config.yml"))
		assert.Equal(t, result.BinPath, path.Join("../testdata/simple/.local/bin"))
		assert.Equal(t, result.Platform, runtime.GOOS)
		assert.Equal(t, result.Architecture, runtime.GOARCH)

		assert.Equal(t, result.Config, &configRoot{
			Packages: map[string]string{
				"fzf": "0.45.0",
			},
		})

		assert.Len(t, result.Binaries, 1)
		assert.Equal(t, result.Binaries[0].Name, "fzf")
	})

	t.Run("ConfigOverride", func(t *testing.T) {
		override := path.Join(emptyConfigDir, ".grab")

		result, err := NewContext(override, "")
		assert.NoError(t, err)

		assert.Equal(t, result.ConfigPath, path.Join(emptyConfigDir, ".grab/config.yml"))
		assert.Equal(t, result.BinPath, path.Join("../testdata/simple/.local/bin"))
		assert.Equal(t, result.Config, &configRoot{
			Packages: map[string]string{},
		})
		assert.Len(t, result.Binaries, 0)
	})

	t.Run("BinOverride", func(t *testing.T) {
		override := t.TempDir()

		result, err := NewContext("", override)
		assert.NoError(t, err)

		assert.Equal(t, result.ConfigPath, path.Join("../testdata/simple/.grab/config.yml"))
		assert.Equal(t, result.BinPath, override)

		assert.Equal(t, result.Config, &configRoot{
			Packages: map[string]string{
				"fzf": "0.45.0",
			},
		})

		assert.Len(t, result.Binaries, 1)
	})
}
