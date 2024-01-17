package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRepositoryValid(t *testing.T) {
	actual, err := loadRepository("testdata/repository/valid")

	expected := &repository{
		Packages: []*configPackage{
			{
				APIVersion: "grab.noizwaves.com/v1alpha1",
				Kind:       "Package",
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
				APIVersion: "grab.noizwaves.com/v1alpha1",
				Kind:       "Package",
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

func TestLoadRepositoryBlank(t *testing.T) {
	_, err := loadRepository("testdata/repository/blank")

	assert.Error(t, err)
	assert.ErrorContains(t, err, "error parsing package YAML")
}

func TestLoadConfigValid(t *testing.T) {
	actual, err := loadConfig("testdata/configs/valid.yml")

	expected := &configRoot{
		Packages: map[string]string{
			"bar": "1.2.0",
			"baz": "0.16.5",
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestLoadConfigBlank(t *testing.T) {
	_, err := loadConfig("testdata/configs/blank.yml")

	assert.Error(t, err)
	assert.ErrorContains(t, err, "error parsing config YAML")
}
