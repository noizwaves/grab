package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryGetUrl(t *testing.T) {
	base := Binary{
		Name:        "foo",
		Version:     "1.2.3",
		TemplateUrl: "https://foo/{{ .Version }}/foo",
	}

	t.Run("Simple", func(t *testing.T) {
		result, err := base.GetUrl("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "https://foo/1.2.3/foo", result)
	})

	t.Run("AllVariables", func(t *testing.T) {
		binary := base
		binary.TemplateUrl = "https://foo/{{ .Version }}/foo-{{ .Platform }}-{{ .Arch }}{{ .Ext }}"

		result, err := binary.GetUrl("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "https://foo/1.2.3/foo-linux-arm64", result)
	})

	t.Run("WithOverrides", func(t *testing.T) {
		binary := base
		binary.TemplateUrl = "https://foo/{{ .Version }}/foo-{{ .Platform }}-{{ .Arch }}{{ .Ext }}"
		binary.Overrides = map[string]Override{
			"linux,arm64": {
				Platform:     "QuantumOS",
				Architecture: "200qbit",
				Extension:    ".zip",
			},
		}

		result, err := binary.GetUrl("linux", "arm64")

		assert.NoError(t, err)
		assert.Equal(t, "https://foo/1.2.3/foo-QuantumOS-200qbit.zip", result)
	})

	t.Run("InvalidTemplate", func(t *testing.T) {
		binary := base
		binary.TemplateUrl = "https://foo/{{ .Version"

		_, err := binary.GetUrl("linux", "arm64")

		assert.ErrorContains(t, err, "Error parsing source template")
	})

	t.Run("InvalidVariable", func(t *testing.T) {
		binary := base
		binary.TemplateUrl = "https://foo/{{ .DoesNotExist }}"

		_, err := binary.GetUrl("linux", "arm64")

		assert.ErrorContains(t, err, "Error rendering source template")
	})
}
