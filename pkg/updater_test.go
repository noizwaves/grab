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

	context, err := NewContext(configDir, t.TempDir())
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
	err = updater.Update(context, out)

	assert.NoError(t, err)
	assert.Contains(t, out.String(), "bar: 1.0.0 -> 2.0.0 (https://fakegithub.com/release-information)")

	asserth.FileContents(t, path.Join(configDir, "config.yml"), "packages:\n  bar: 2.0.0\n")
}
