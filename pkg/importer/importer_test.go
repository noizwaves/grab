package importer

import (
	"bytes"
	"testing"

	"github.com/noizwaves/grab/pkg"
	"github.com/noizwaves/grab/pkg/github"
	"github.com/noizwaves/grab/pkg/internal/githubh"
	"github.com/noizwaves/grab/pkg/internal/osh"
	"github.com/stretchr/testify/assert"
)

// TestImport_Success_LatestRelease tests importing from a latest release URL.
func TestImport_Success_LatestRelease(t *testing.T) {
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
	err = importer.Import(gCtx, "https://github.com/foo/myapp/releases/latest", &out)

	assert.NoError(t, err)
	assert.Contains(t, out.String(), `Package "myapp" saved to`)
	assert.Contains(t, out.String(), `/myapp.yml`)
}

// TestImport_Success_TaggedRelease tests importing from a tagged release URL.
func TestImport_Success_TaggedRelease(t *testing.T) {
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
	err = importer.Import(gCtx, "https://github.com/foo/myapp/releases/tag/v1.2.3", &out)

	assert.NoError(t, err)
	assert.Contains(t, out.String(), `Package "myapp" saved to`)
	assert.Contains(t, out.String(), `/myapp.yml`)
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
