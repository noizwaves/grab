package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRepositoryValid(t *testing.T) {
	actual, err := loadRepository("testdata/configs/valid")

	expected := repository{
		Packages: []configPackage{
			{
				Metadata: configPackageMetadata{
					Name: "bar",
				},
				Spec: configPackageSpec{
					GitHubRelease: configGitHubRelease{
						Org:          "foo",
						Repo:         "bar",
						Name:         "{{ .Version }}",
						VersionRegex: "\\d+\\.\\d+\\.\\d+",
						FileName: map[string]string{
							"darwin,amd64": "bin",
							"darwin,arm64": "bin",
							"linux,amd64":  "bin",
							"linux,arm64":  "bin",
						},
					},
					Program: configProgram{
						VersionArgs:  []string{"--version"},
						VersionRegex: "\\d+\\.\\d+\\.\\d+",
					},
				},
			},
			{
				Metadata: configPackageMetadata{
					Name: "baz",
				},
				Spec: configPackageSpec{
					GitHubRelease: configGitHubRelease{
						Org:          "foo",
						Repo:         "baz",
						Name:         "v{{ .Version }}",
						VersionRegex: "\\d+\\.\\d+\\.\\d+",
						FileName: map[string]string{
							"darwin,amd64": "v{{ .Version }}-apple-darwin-x86_64.zip",
							"darwin,arm64": "v{{ .Version }}-apple-darwin-aarch64.zip",
							"linux,amd64":  "v{{ .Version }}-unknown-linux-musl-x86_64.tgz",
							"linux,arm64":  "v{{ .Version }}-unknown-linux-gnu-aarch64.tar.gz",
						},
					},
					Program: configProgram{
						VersionArgs:  []string{"version"},
						VersionRegex: "\\d+\\.\\d+\\.\\d+",
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestLoadConfigValid(t *testing.T) {
	actual, err := loadConfig("testdata/configs/valid")

	expected := configRoot{
		Packages: map[string]string{
			"bar": "1.2.0",
			"baz": "0.16.5",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
