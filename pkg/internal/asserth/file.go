package asserth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FileContents(t *testing.T, path, expected string) {
	t.Helper()

	actual, err := os.ReadFile(path)
	assert.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}
