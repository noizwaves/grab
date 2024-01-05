package pkg

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderSourceUrl(t *testing.T) {
	base := configBinary{
		Name:    "foo",
		Version: "1.2.3",
		Source:  "https://foo/{{ .Version }}/foo",
	}

	t.Run("Simple", func(t *testing.T) {
		result, err := renderSourceUrl(base)

		assert.NoError(t, err)
		assert.Equal(t, "https://foo/1.2.3/foo", result)
	})

	t.Run("AllVariables", func(t *testing.T) {
		binary := base
		binary.Source = "https://foo/{{ .Version }}/foo-{{ .Platform }}-{{ .Arch }}"

		result, err := renderSourceUrl(binary)

		assert.NoError(t, err)
		expected := fmt.Sprintf("https://foo/1.2.3/foo-%s-%s", runtime.GOOS, runtime.GOARCH)
		assert.Equal(t, expected, result)
	})

	t.Run("WithOverrides", func(t *testing.T) {
		binary := base
		binary.Source = "https://foo/{{ .Version }}/foo-{{ .Platform }}-{{ .Arch }}"
		binary.Platforms = map[string]map[string][]string{
			runtime.GOOS: {
				runtime.GOARCH: {"QuantumOS", "200qbit"},
			},
		}

		result, err := renderSourceUrl(binary)

		assert.NoError(t, err)
		expected := "https://foo/1.2.3/foo-QuantumOS-200qbit"
		assert.Equal(t, expected, result)
	})

	t.Run("InvalidTemplate", func(t *testing.T) {
		binary := base
		binary.Source = "https://foo/{{ .Version"

		_, err := renderSourceUrl(binary)

		assert.ErrorContains(t, err, "Error parsing Source as template")
	})

	t.Run("InvalidVariable", func(t *testing.T) {
		binary := base
		binary.Source = "https://foo/{{ .DoesNotExist }}"

		_, err := renderSourceUrl(binary)

		assert.ErrorContains(t, err, "Error rendering Source as template")
	})
}
