package pkg

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/noizwaves/grab/pkg/internal/asserth"
	"github.com/noizwaves/grab/pkg/internal/githubh"
	"github.com/noizwaves/grab/pkg/internal/osh"
	"github.com/stretchr/testify/assert"
)

// Simple test case that installs one package into an empty bin directory.
func TestInstall(t *testing.T) {
	configDir := osh.CopyDir(t, "testdata/contexts/simple")
	binDir := t.TempDir()

	context, err := NewContext(configDir, binDir)
	if err != nil {
		t.Fatal(err)
	}

	installer := Installer{
		GitHubClient: &githubh.MockGitHubClient{
			AssetData: []byte("#!/usr/bin/env bash\necho '1.0.0'"),
		},
	}

	out := bytes.Buffer{}
	err = installer.Install(context, &out)

	assert.NoError(t, err)
	assert.Contains(t, out.String(), "bar: installing 1.0.0... Done!")

	barPath := filepath.Join(binDir, "bar")
	assert.FileExists(t, barPath)
	asserth.CommandSucceeds(t, barPath)
	asserth.CommandStdoutContains(t, barPath, "1.0.0")
}
