package pkg

import (
	"bytes"
	"path"
	"testing"

	"github.com/noizwaves/grab/pkg/github"
	"github.com/noizwaves/grab/pkg/internal/asserth"
	"github.com/noizwaves/grab/pkg/internal/githubh"
	"github.com/noizwaves/grab/pkg/internal/osh"
	"github.com/stretchr/testify/assert"
)

// Simple test that updates the version of a single package.
func TestUpdate(t *testing.T) {
	configDir := osh.CopyDir(t, "testdata/contexts/simple")

	gCtx, err := NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	updater := Updater{
		GitHubClient: &githubh.MockGitHubClient{
			Release: &github.Release{
				Name: "2.0.0",
				URL:  "https://fakegithub.com/release-information",
			},
		},
	}

	out := &bytes.Buffer{}
	err = updater.Update(gCtx, "", out)

	assert.NoError(t, err)

	// Verify consistent output formatting
	output := out.String()
	assert.Contains(t, output, "bar: 1.0.0 -> 2.0.0 (https://fakegithub.com/release-information)")
	assert.Contains(t, output, "Updated config file. Now run `grab install`.")

	asserth.FileContents(t, path.Join(configDir, "config.yml"), "packages:\n  bar: 2.0.0\n")
}

// Test that validates a package name exists in configuration.
func TestUpdateValidPackageName(t *testing.T) {
	configDir := osh.CopyDir(t, "testdata/contexts/simple")

	gCtx, err := NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	updater := Updater{
		GitHubClient: &githubh.MockGitHubClient{
			Release: &github.Release{
				Name: "2.0.0",
				URL:  "https://fakegithub.com/release-information",
			},
		},
	}

	out := &bytes.Buffer{}
	err = updater.Update(gCtx, "bar", out)

	assert.NoError(t, err)

	// Verify consistent output formatting for single package update
	output := out.String()
	assert.Contains(t, output, "bar: 1.0.0 -> 2.0.0 (https://fakegithub.com/release-information)")
	assert.Contains(t, output, "Updated config file. Now run `grab install`.")
}

// Test that returns error when package name doesn't exist in configuration.
func TestUpdateInvalidPackageName(t *testing.T) {
	configDir := osh.CopyDir(t, "testdata/contexts/simple")

	gCtx, err := NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	updater := Updater{
		GitHubClient: &githubh.MockGitHubClient{
			Release: &github.Release{
				Name: "2.0.0",
				URL:  "https://fakegithub.com/release-information",
			},
		},
	}

	out := &bytes.Buffer{}
	err = updater.Update(gCtx, "nonexistent", out)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), `package "nonexistent" not found in configuration`)
}

// Test that single package update only processes the specified package and ignores others.
func TestUpdateSinglePackageIgnoresOthers(t *testing.T) {
	configDir := osh.CopyDir(t, "testdata/contexts/multiple")

	gCtx, err := NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	mockClient := &githubh.MockGitHubClient{
		Release: &github.Release{
			Name: "2.0.0",
			URL:  "https://fakegithub.com/release-information",
		},
	}

	updater := Updater{
		GitHubClient: mockClient,
	}

	out := &bytes.Buffer{}
	err = updater.Update(gCtx, "bar", out)

	assert.NoError(t, err)

	// Verify only one API call was made
	assert.Equal(t, 1, len(mockClient.GetLatestReleaseCalls), "Expected exactly 1 API call")
	assert.Equal(t, "bar", mockClient.GetLatestReleaseCalls[0].Repo, "Expected only 'bar' package to be requested")

	// Verify output only contains "bar" and not "baz"
	output := out.String()
	assert.Contains(t, output, "bar: 1.0.0 -> 2.0.0")
	assert.NotContains(t, output, "baz", "Output should not contain 'baz' package")

	// Verify that the config file shows only "bar" was updated, "baz" unchanged
	asserth.FileContents(t, path.Join(configDir, "config.yml"), "packages:\n  bar: 2.0.0\n  baz: 1.2.3\n")
}

// Test that "is latest" output format is consistent for single package updates.
func TestUpdatePackageAlreadyLatest(t *testing.T) {
	configDir := osh.CopyDir(t, "testdata/contexts/simple")

	gCtx, err := NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	updater := Updater{
		GitHubClient: &githubh.MockGitHubClient{
			Release: &github.Release{
				Name: "1.0.0", // Same version as in config
				URL:  "https://fakegithub.com/release-information",
			},
		},
	}

	out := &bytes.Buffer{}
	err = updater.Update(gCtx, "bar", out)

	assert.NoError(t, err)

	// Verify "is latest" output format for single package
	output := out.String()
	assert.Contains(t, output, "bar: 1.0.0 is latest")

	// Should not contain config update message when no changes
	assert.NotContains(t, output, "Updated config file")
}

// Test error handling for no packages configured.
func TestUpdateNoPackagesConfigured(t *testing.T) {
	configDir := osh.CopyDir(t, "testdata/contexts/empty")

	gCtx, err := NewGrabContext(configDir, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	updater := Updater{
		GitHubClient: &githubh.MockGitHubClient{
			Release: &github.Release{
				Name: "2.0.0",
				URL:  "https://fakegithub.com/release-information",
			},
		},
	}

	out := &bytes.Buffer{}
	err = updater.Update(gCtx, "", out)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no packages configured")
}
