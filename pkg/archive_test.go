package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnTgzFileNamedValid(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		BinaryName     string
		ExpectedErr    string
		ExpectedResult string
	}{
		{
			Name:           "Matches",
			BinaryName:     "binary",
			ExpectedResult: "foobar\n",
		},
		{
			Name:        "Different",
			BinaryName:  "other",
			ExpectedErr: "No file named \"other\" found in archive",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			f, err := os.Open("testdata/archives/binary.tgz")
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			result, err := unTgzFileNamed(tc.BinaryName, f)

			if tc.ExpectedResult != "" {
				assert.Equal(t, []byte(tc.ExpectedResult), result)
			} else {
				assert.ErrorContains(t, err, tc.ExpectedErr)
			}
		})
	}
}

func TestUnTgzFileNamedInvalid(t *testing.T) {
	f, err := os.Open("testdata/archives/binary")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = unTgzFileNamed("binary", f)

	assert.ErrorContains(t, err, "Error decompressing Gzipped data")
}

func TestUnGzipValid(t *testing.T) {
	f, err := os.Open("testdata/archives/binary.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	result, err := unGzip(f)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte("foobar\n"), result)
}

func TestUnGzipInvalid(t *testing.T) {
	f, err := os.Open("testdata/archives/binary")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = unGzip(f)

	assert.ErrorContains(t, err, "Error decompressing Gzipped data")
}
