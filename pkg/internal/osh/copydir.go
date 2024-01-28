package osh

import (
	"os"
	"path/filepath"
	"testing"
)

// Copy an existing directory to a temporary directory for testing.
func CopyDir(t *testing.T, source string) string {
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
