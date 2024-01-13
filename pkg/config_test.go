package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigValid(t *testing.T) {
	actual, err := loadConfig("testdata/configs/valid.yml")

	expected := configRoot{
		Binaries: []configBinary{
			{
				Name:    "bar",
				Version: "1.2.0",
				Source: configSource{
					Org:          "foo",
					Repo:         "bar",
					ReleaseName:  "{{ .Version }}",
					ReleaseRegex: ".*",
					FileName:     "bin",
					VersionFlags: []string{"--version"},
					VersionRegex: "\\d+\\.\\d+\\.\\d+",
				},
			},
			{
				Name:    "baz",
				Version: "0.16.5",
				Source: configSource{
					Org:          "foo",
					Repo:         "baz",
					ReleaseName:  "v{{ .Version }}",
					ReleaseRegex: "v.*",
					FileName:     "v{{ .Version }}-{{ .Platform }}-{{ .Arch }}.{{ .Ext }}",
					Overrides: map[configPlatformKey]configPlatformValue{
						"linux": {
							"amd64": [3]string{"unknown-linux-musl", "x86_64", "tgz"},
							"arm64": [3]string{"unknown-linux-gnu", "aarch64", "tar.gz"},
						},
						"darwin": {
							"amd64": [3]string{"apple-darwin", "x86_64", "zip"},
							"arm64": [3]string{"apple-darwin", "aarch64", "zip"},
						},
					},
					VersionFlags: []string{"version"},
					VersionRegex: "\\d+\\.\\d+\\.\\d+",
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
