package importer

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/noizwaves/grab/pkg"
	"github.com/noizwaves/grab/pkg/github"
	"github.com/noizwaves/grab/pkg/internal/githubh"
	"github.com/noizwaves/grab/pkg/internal/osh"
	"github.com/stretchr/testify/assert"
)

// Happy path test.
func TestImport_Success(t *testing.T) {
	configDir := osh.CopyDir(t, "../testdata/contexts/simple")

	gCtx, err := pkg.NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	mockRelease := &github.Release{
		TagName: "v1.2.3",
		Assets: []github.Asset{
			{Name: "myapp-1.2.3-linux-amd64.tar.gz"},
			{Name: "myapp-1.2.3-linux-arm64.tar.gz"},
			{Name: "myapp-1.2.3-darwin-amd64.tar.gz"},
			{Name: "myapp-1.2.3-darwin-arm64.tar.gz"},
		},
	}

	importer := NewImporter(&githubh.MockGitHubClient{
		Release: mockRelease,
	})

	out := bytes.Buffer{}
	err = importer.Import(gCtx, "https://github.com/foo/myapp", &out)

	assert.NoError(t, err)
	assert.Contains(t, out.String(), `Package "myapp" saved to`)
	assert.Contains(t, out.String(), `/myapp.yml`)

	// Verify the contents of the package YAML file
	expectedYamlPath := filepath.Join(configDir, "repository", "myapp.yml")
	actualYaml, err := os.ReadFile(expectedYamlPath)
	assert.NoError(t, err)

	expectedYaml := `apiVersion: grab.noizwaves.com/v1alpha1
kind: Package
metadata:
  name: myapp
spec:
  gitHubRelease:
    org: foo
    repo: myapp
    name: v{{ .Version }}
    versionRegex: \d+\.\d+\.\d+
    fileName:
      darwin,amd64: myapp-{{ .Version }}-darwin-amd64.tar.gz
      darwin,arm64: myapp-{{ .Version }}-darwin-arm64.tar.gz
      linux,amd64: myapp-{{ .Version }}-linux-amd64.tar.gz
      linux,arm64: myapp-{{ .Version }}-linux-arm64.tar.gz
  program:
    versionArgs: [--version]
    versionRegex: \d+\.\d+\.\d+
`

	assert.Equal(t, expectedYaml, string(actualYaml))
}

// TestImport_Error_InvalidURL tests error handling for invalid URLs.
func TestImport_Error_InvalidURL(t *testing.T) {
	configDir := osh.CopyDir(t, "../testdata/contexts/simple")

	gCtx, err := pkg.NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	importer := NewImporter(&githubh.MockGitHubClient{})

	out := bytes.Buffer{}
	err = importer.Import(gCtx, "https://example.com/invalid/url", &out)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

// TestImport_Error_GitHubAPIFailure tests error handling when GitHub API fails.
func TestImport_Error_GitHubAPIFailure(t *testing.T) {
	configDir := osh.CopyDir(t, "../testdata/contexts/simple")

	gCtx, err := pkg.NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	// MockGitHubClient with no Release set will return "not implemented" error
	importer := NewImporter(&githubh.MockGitHubClient{})

	out := bytes.Buffer{}
	err = importer.Import(gCtx, "https://github.com/foo/myapp/releases/latest", &out)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get release")
}

// TestImport_Error_NoMatchingAssets tests error handling when no assets match platform/arch.
func TestImport_Error_NoMatchingAssets(t *testing.T) {
	configDir := osh.CopyDir(t, "../testdata/contexts/simple")

	gCtx, err := pkg.NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	// Release with assets that don't match any platform/arch patterns
	mockRelease := &github.Release{
		TagName: "v1.2.3",
		Assets: []github.Asset{
			{Name: "myapp-1.2.3-unknown-platform.tar.gz"},
			{Name: "myapp-1.2.3-windows-x72.zip"},
		},
	}

	importer := NewImporter(&githubh.MockGitHubClient{
		Release: mockRelease,
	})

	out := bytes.Buffer{}
	err = importer.Import(gCtx, "https://github.com/foo/myapp/releases/latest", &out)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no matching asset name found")
}
