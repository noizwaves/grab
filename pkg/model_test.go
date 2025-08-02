package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAssetFileName(t *testing.T) {
	base := Binary{
		Name:          "foo",
		PinnedVersion: "1.2.3",
		Org:           "bar",
		Repo:          "foo",
		releaseName:   "{{ .Version }}",
		fileName: map[string]string{
			"linux,arm64": "foo",
		},
	}

	t.Run("NoVariables", func(t *testing.T) {
		result, err := base.GetAssetFileName("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "foo", result)
	})

	t.Run("AllVariables", func(t *testing.T) {
		binary := base
		binary.fileName = map[string]string{
			"linux,arm64": "foo-{{ .Version }}",
		}

		result, err := binary.GetAssetFileName("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "foo-1.2.3", result)
	})

	t.Run("InvalidFileNameTemplate", func(t *testing.T) {
		binary := base
		binary.fileName = map[string]string{
			"linux,arm64": "foo-{{ .Version",
		}

		_, err := binary.GetAssetFileName("linux", "arm64")

		assert.ErrorContains(t, err, "error parsing asset filename template")
	})

	t.Run("InvalidVariable", func(t *testing.T) {
		binary := base
		binary.fileName = map[string]string{
			"linux,arm64": "foo-{{ .DoesNotExist }}",
		}

		_, err := binary.GetAssetFileName("linux", "arm64")

		assert.ErrorContains(t, err, "error rendering asset filename template")
	})
}

func TestGetReleaseName(t *testing.T) {
	base := Binary{
		Name:          "foo",
		PinnedVersion: "1.2.3",
		Org:           "bar",
		Repo:          "foo",
		releaseName:   "{{ .Version }}",
		fileName: map[string]string{
			"linux,arm64": "foo",
		},
	}

	t.Run("Simple", func(t *testing.T) {
		result, err := base.GetReleaseName()

		assert.NoError(t, err)
		assert.Equal(t, "1.2.3", result)
	})

	t.Run("InvalidReleaseNameTemplate", func(t *testing.T) {
		binary := base
		binary.releaseName = "v{{ .Version"

		_, err := binary.GetReleaseName()

		assert.ErrorContains(t, err, "error parsing release name template")
	})

	t.Run("InvalidVariable", func(t *testing.T) {
		binary := base
		binary.releaseName = "v-{{ .DoesNotExist }}"

		_, err := binary.GetReleaseName()

		assert.ErrorContains(t, err, "error rendering release name template")
	})
}

func TestGetEmbeddedBinaryPath(t *testing.T) {
	base := Binary{
		Name:          "foo",
		PinnedVersion: "1.2.3",
		Org:           "bar",
		Repo:          "foo",
		releaseName:   "{{ .Version }}",
		fileName: map[string]string{
			"linux,arm64": "{{ .Version }}",
		},
	}

	t.Run("DefaultsToPackageName", func(t *testing.T) {
		result, err := base.GetEmbeddedBinaryPath("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "foo", result)
	})

	t.Run("MatchingEmbeddedBinaryPath", func(t *testing.T) {
		binary := base
		binary.embeddedBinaryPath = map[string]string{
			"linux,arm64": "linux-arm64/foo",
		}

		result, err := binary.GetEmbeddedBinaryPath("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "linux-arm64/foo", result)
	})

	t.Run("MissingEmbeddedBinaryPath", func(t *testing.T) {
		binary := base
		binary.embeddedBinaryPath = map[string]string{}

		result, err := binary.GetEmbeddedBinaryPath("linux", "amd64")

		assert.Equal(t, "", result)
		assert.ErrorContains(t, err, "missing value for platform=linux,arch=amd64")
	})

	t.Run("EmbeddedBinaryPathWithVersion", func(t *testing.T) {
		binary := base
		binary.embeddedBinaryPath = map[string]string{
			"linux,arm64": "bin/{{ .Version }}/foo",
		}

		result, err := binary.GetEmbeddedBinaryPath("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "bin/1.2.3/foo", result)
	})

	t.Run("InvalidEmbeddedBinaryPathTemplate", func(t *testing.T) {
		binary := base
		binary.embeddedBinaryPath = map[string]string{
			"linux,arm64": "bin/{{ .Version",
		}

		_, err := binary.GetEmbeddedBinaryPath("linux", "arm64")

		assert.ErrorContains(t, err, "error parsing embedded binary path template")
	})

	t.Run("InvalidEmbeddedBinaryPathVariable", func(t *testing.T) {
		binary := base
		binary.embeddedBinaryPath = map[string]string{
			"linux,arm64": "bin/{{ .DoesNotExist }}/foo",
		}

		_, err := binary.GetEmbeddedBinaryPath("linux", "arm64")

		assert.ErrorContains(t, err, "error rendering embedded binary path template")
	})
}

func TestBinaryShouldReplace(t *testing.T) {
	base := Binary{
		Name:          "foo",
		PinnedVersion: "1.2.3",
		Org:           "bar",
		Repo:          "foo",
		releaseName:   "{{ .Version }}",
		fileName: map[string]string{
			"linux,arm64": "foo",
		},
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
