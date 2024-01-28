package pkg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/noizwaves/grab/pkg/github"
	"github.com/stretchr/testify/assert"
)

type mockGitHubClient struct {
	assetData []byte
	release   *github.Release
}

func (m *mockGitHubClient) DownloadReleaseAsset(_, _, _, _ string) ([]byte, error) {
	if len(m.assetData) == 0 {
		return nil, fmt.Errorf("not implemented")
	}

	return m.assetData, nil
}

func (m *mockGitHubClient) GetLatestRelease(_, _ string) (*github.Release, error) {
	if m.release == nil {
		return nil, fmt.Errorf("not implemented")
	}

	return m.release, nil
}

func assertCommandSucceeds(t *testing.T, path string) {
	t.Helper()

	cmd := exec.Command(path)
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = cmd.Wait()

	if err != nil {
		assert.Fail(t, "command did not run successfully", err)
	}
}

func assertCommandStdoutContains(t *testing.T, path string, expected string) {
	t.Helper()

	cmd := exec.Command(path)

	out := bytes.Buffer{}
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	stdout := out.String()
	assert.Contains(t, stdout, expected)
}

// Simple test case that installs one package into an empty bin directory.
func TestInstall(t *testing.T) {
	configDir := copyExistingDir(t, "testdata/contexts/simple")
	binDir := t.TempDir()

	context, err := NewContext(configDir, binDir)
	if err != nil {
		t.Fatal(err)
	}

	installer := Installer{
		GitHubClient: &mockGitHubClient{
			assetData: []byte("#!/usr/bin/env bash\necho '1.0.0'"),
		},
	}

	out := bytes.Buffer{}
	err = installer.Install(context, &out)

	assert.NoError(t, err)
	assert.Contains(t, out.String(), "bar: installing 1.0.0... Done!")

	barPath := filepath.Join(binDir, "bar")
	assert.FileExists(t, barPath)
	assertCommandSucceeds(t, barPath)
	assertCommandStdoutContains(t, barPath, "1.0.0")
}

func copyExistingDir(t *testing.T, source string) string {
	t.Helper()

	destination := t.TempDir()

	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatal(err)
		}

		relativePath, err := filepath.Rel(source, path)
		if err != nil {
			t.Fatal(err)
		}

		if info.IsDir() {
			err = os.MkdirAll(filepath.Join(destination, relativePath), info.Mode())
			if err != nil {
				t.Fatal(err)
			}
		} else {
			contents, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			err = os.WriteFile(filepath.Join(destination, relativePath), contents, info.Mode())
			if err != nil {
				t.Fatal(err)
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return destination
}
