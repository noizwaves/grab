package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryGetURL(t *testing.T) { //nolint:funlen
	base := Binary{
		Name:        "foo",
		Version:     "1.2.3",
		Org:         "bar",
		Repo:        "foo",
		ReleaseName: "{{ .Version }}",
		FileName:    "foo",
	}

	t.Run("Simple", func(t *testing.T) {
		result, err := base.GetURL("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "https://github.com/bar/foo/releases/download/1.2.3/foo", result)
	})

	t.Run("AllVariables", func(t *testing.T) {
		binary := base
		binary.FileName = "foo-{{ .Platform }}-{{ .Arch }}{{ .Ext }}"

		result, err := binary.GetURL("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "https://github.com/bar/foo/releases/download/1.2.3/foo-linux-arm64", result)
	})

	t.Run("WithOverrides", func(t *testing.T) {
		binary := base
		binary.FileName = "foo-{{ .Platform }}-{{ .Arch }}{{ .Ext }}"
		binary.Overrides = map[string]Override{
			"linux,arm64": {
				Platform:     "QuantumOS",
				Architecture: "200qbit",
				Extension:    ".zip",
			},
		}

		result, err := binary.GetURL("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "https://github.com/bar/foo/releases/download/1.2.3/foo-QuantumOS-200qbit.zip", result)
	})

	t.Run("InvalidReleaseNameTemplate", func(t *testing.T) {
		binary := base
		binary.ReleaseName = "v{{ .Version"

		_, err := binary.GetURL("linux", "arm64")

		assert.ErrorContains(t, err, "error parsing source template")
	})

	t.Run("InvalidFileNameTemplate", func(t *testing.T) {
		binary := base
		binary.FileName = "foo-{{ .Version"

		_, err := binary.GetURL("linux", "arm64")

		assert.ErrorContains(t, err, "error parsing source template")
	})

	t.Run("InvalidVariable", func(t *testing.T) {
		binary := base
		binary.ReleaseName = "v-{{ .DoesNotExist }}"

		_, err := binary.GetURL("linux", "arm64")

		assert.ErrorContains(t, err, "error rendering source template")
	})
}

func TestBinaryShouldReplace(t *testing.T) {
	base := Binary{
		Name:        "foo",
		Version:     "1.2.3",
		Org:         "bar",
		Repo:        "foo",
		ReleaseName: "{{ .Version }}",
		FileName:    "foo",
	}

	t.Run("CurrentLessThanDesired", func(t *testing.T) {
		assert.True(t, base.ShouldReplace("1.0.0"))
	})

	t.Run("CurrentGreaterThanDesired", func(t *testing.T) {
		assert.True(t, base.ShouldReplace("9.9.9"))
	})

	t.Run("SameValue", func(t *testing.T) {
		assert.False(t, base.ShouldReplace("1.2.3"))
	})
}
