package asserth

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CommandSucceeds(t *testing.T, path string) {
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

func CommandStdoutContains(t *testing.T, path string, expected string) {
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
