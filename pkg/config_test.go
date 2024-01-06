package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigValid(t *testing.T) {
	actual, err := parseConfig("testdata/configs/valid.yml")

	expected := configRoot{
		Binaries: []configBinary{
			{
				Name:    "foo",
				Source:  "https://foo.com/{{ .Version }}/bin",
				Version: "1.2.0",
			},
			{
				Name:    "bar",
				Source:  "https://bar.com/releases/v{{ .Version }}/assets/v{{ .Version }}-{{ .Platform }}-{{ .Arch }}.{{ .Ext }}",
				Version: "0.16.5",
				Platforms: map[configPlatformKey]configPlatformValue{
					"linux": {
						"amd64": []string{"unknown-linux-musl", "x86_64", "tgz"},
						"arm64": []string{"unknown-linux-gnu", "aarch64", "tar.gz"},
					},
					"darwin": {
						"amd64": []string{"apple-darwin", "x86_64", "zip"},
						"arm64": []string{"apple-darwin", "aarch64", "zip"},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
