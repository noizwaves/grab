package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnTgzFileNamedValid(t *testing.T) {
	for _, testCase := range []struct {
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
			ExpectedErr: "no file named \"other\" found in archive",
		},
	} {
		t.Run(testCase.Name, func(t *testing.T) {
			file, err := os.Open("testdata/archives/binary.tgz")
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			result, err := unTgzFileNamed(testCase.BinaryName, file)

			if testCase.ExpectedResult != "" {
				assert.Equal(t, []byte(testCase.ExpectedResult), result)
			} else {
				assert.ErrorContains(t, err, testCase.ExpectedErr)
			}
		})
	}
}

func TestUnTgzFileNamedInvalid(t *testing.T) {
	file, err := os.Open("testdata/archives/binary")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = unTgzFileNamed("binary", file)

	assert.ErrorContains(t, err, "error decompressing Gzipped data")
}

func TestUnZipFileNamedValid(t *testing.T) {
	for _, testCase := range []struct {
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
			ExpectedErr: "no file named \"other\" found in archive",
		},
	} {
		t.Run(testCase.Name, func(t *testing.T) {
			file, err := os.Open("testdata/archives/binary.zip")
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			result, err := unZipFileNamed(testCase.BinaryName, file)

			if testCase.ExpectedResult != "" {
				assert.NoError(t, err)
				assert.Equal(t, []byte(testCase.ExpectedResult), result)
			} else {
				assert.ErrorContains(t, err, testCase.ExpectedErr)
			}
		})
	}
}

func TestUnZipFileNamedInvalid(t *testing.T) {
	file, err := os.Open("testdata/archives/binary")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = unZipFileNamed("binary", file)

	assert.ErrorContains(t, err, "error decompressing Zipped data")
}

func TestUnGzipValid(t *testing.T) {
	file, err := os.Open("testdata/archives/binary.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	result, err := unGzip(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte("foobar\n"), result)
}

func TestUnGzipInvalid(t *testing.T) {
	file, err := os.Open("testdata/archives/binary")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = unGzip(file)

	assert.ErrorContains(t, err, "error decompressing Gzipped data")
}
