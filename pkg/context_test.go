package pkg

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/noizwaves/grab/pkg/internal/asserth"
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

	err = os.WriteFile(path.Join(dir, ".grab", "config.yml"), emptyConfig, 0o644)
	if err != nil {
		t.Fatal(err)
	}

	return dir
}

func TestNewGrabContext(t *testing.T) {
	// Use simple as $HOME
	t.Setenv("HOME", "../testdata/simple")

	// Set up an overrides directory
	emptyConfigDir := makeEmptyConfig(t)

	t.Run("NoOverrides", func(t *testing.T) {
		result, err := NewGrabContext("", "")

		assert.NoError(t, err)

		assert.Equal(t, result.ConfigPath, path.Join("../testdata/simple/.grab/config.yml"))
		assert.Equal(t, result.BinPath, path.Join("../testdata/simple/.local/bin"))
		assert.Equal(t, result.RepoPath, path.Join("../testdata/simple/.grab/repository"))
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

		result, err := NewGrabContext(override, "")
		assert.NoError(t, err)

		assert.Equal(t, result.ConfigPath, path.Join(emptyConfigDir, ".grab/config.yml"))
		assert.Equal(t, result.RepoPath, path.Join(emptyConfigDir, ".grab/repository"))
		assert.Equal(t, result.BinPath, path.Join("../testdata/simple/.local/bin"))
		assert.Equal(t, result.Config, &configRoot{
			Packages: map[string]string{},
		})
		assert.Len(t, result.Binaries, 0)
	})

	t.Run("BinOverride", func(t *testing.T) {
		override := t.TempDir()

		result, err := NewGrabContext("", override)
		assert.NoError(t, err)

		assert.Equal(t, result.ConfigPath, path.Join("../testdata/simple/.grab/config.yml"))
		assert.Equal(t, result.RepoPath, path.Join("../testdata/simple/.grab/repository"))
		assert.Equal(t, result.BinPath, override)

		assert.Equal(t, result.Config, &configRoot{
			Packages: map[string]string{
				"fzf": "0.45.0",
			},
		})

		assert.Len(t, result.Binaries, 1)
	})
}

func TestAddPackageToConfig(t *testing.T) {
	dir := makeEmptyConfig(t)
	configDirPath := path.Join(dir, ".grab")

	gCtx, err := NewGrabContext(configDirPath, t.TempDir())
	assert.NoError(t, err)

	err = gCtx.AddPackageToConfig("fzf", "0.45.0")
	assert.NoError(t, err)

	// Verify in-memory state
	assert.Equal(t, "0.45.0", gCtx.Config.Packages["fzf"])

	// Verify persisted to disk
	expectedContent := "packages:\n" +
		"  fzf: 0.45.0\n"
	asserth.FileContents(t, path.Join(configDirPath, "config.yml"), expectedContent)
}

func TestAddPackageToConfigAddsToExisting(t *testing.T) {
	dir := makeEmptyConfig(t)
	configDirPath := path.Join(dir, ".grab")

	gCtx, err := NewGrabContext(configDirPath, t.TempDir())
	assert.NoError(t, err)

	// Add first package
	err = gCtx.AddPackageToConfig("bar", "1.0.0")
	assert.NoError(t, err)

	// Add second package
	err = gCtx.AddPackageToConfig("fzf", "0.45.0")
	assert.NoError(t, err)

	// Verify both packages are in memory
	assert.Equal(t, "1.0.0", gCtx.Config.Packages["bar"])
	assert.Equal(t, "0.45.0", gCtx.Config.Packages["fzf"])

	// Verify persisted to disk
	expectedContent := "packages:\n" +
		"  bar: 1.0.0\n" +
		"  fzf: 0.45.0\n"
	asserth.FileContents(t, path.Join(configDirPath, "config.yml"), expectedContent)
}
