package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRepositoryValid(t *testing.T) {
	actual, err := loadRepository("testdata/configs/valid/.garb/repository")

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
							"darwin,amd64": "v{{ .Version }}-x86_64-apple-darwin.zip",
							"darwin,arm64": "v{{ .Version }}-aarch64-apple-darwin.zip",
							"linux,amd64":  "v{{ .Version }}-x86_64-unknown-linux-musl.tgz",
							"linux,arm64":  "v{{ .Version }}-aarch64-unknown-linux-gnu.tar.gz",
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
	actual, err := loadConfig("testdata/configs/valid/.garb/config.yml")

	expected := configRoot{
		Packages: map[string]string{
			"bar": "1.2.0",
			"baz": "0.16.5",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
