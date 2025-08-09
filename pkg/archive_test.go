package pkg

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ulikunitz/xz"
)

func TestListTgzContents(t *testing.T) {
	// Create test tar.gz content
	var buf bytes.Buffer

	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	// Add test files
	files := []struct {
		name    string
		content string
	}{
		{"hyperfine/hyperfine", "binary content"},
		{"hyperfine/README.md", "readme content"},
		{"hyperfine/LICENSE", "license content"},
	}

	for _, file := range files {
		header := &tar.Header{
			Name: file.name,
			Mode: 0o644,
			Size: int64(len(file.content)),
		}
		err := tarWriter.WriteHeader(header)
		require.NoError(t, err)

		_, err = tarWriter.Write([]byte(file.content))
		require.NoError(t, err)
	}

	err := tarWriter.Close()
	require.NoError(t, err)
	err = gzWriter.Close()
	require.NoError(t, err)

	// Test listing contents
	result, err := ListTgzContents(&buf)
	require.NoError(t, err)

	expected := []string{"hyperfine/hyperfine", "hyperfine/README.md", "hyperfine/LICENSE"}
	assert.Equal(t, expected, result)
}

func TestListZipContents(t *testing.T) {
	// Create test zip content
	var buf bytes.Buffer

	zipWriter := zip.NewWriter(&buf)

	// Add test files
	files := []struct {
		name    string
		content string
	}{
		{"hyperfine/hyperfine", "binary content"},
		{"hyperfine/README.md", "readme content"},
		{"hyperfine/LICENSE", "license content"},
	}

	for _, file := range files {
		writer, err := zipWriter.Create(file.name)
		require.NoError(t, err)

		_, err = writer.Write([]byte(file.content))
		require.NoError(t, err)
	}

	err := zipWriter.Close()
	require.NoError(t, err)

	// Test listing contents
	result, err := ListZipContents(&buf)
	require.NoError(t, err)

	expected := []string{"hyperfine/hyperfine", "hyperfine/README.md", "hyperfine/LICENSE"}
	assert.Equal(t, expected, result)
}

func TestListTarxzContents(t *testing.T) {
	// Create test tar content first
	var tarBuf bytes.Buffer

	tarWriter := tar.NewWriter(&tarBuf)

	// Add test files
	files := []struct {
		name    string
		content string
	}{
		{"hyperfine/hyperfine", "binary content"},
		{"hyperfine/README.md", "readme content"},
	}

	for _, file := range files {
		header := &tar.Header{
			Name: file.name,
			Mode: 0o644,
			Size: int64(len(file.content)),
		}
		err := tarWriter.WriteHeader(header)
		require.NoError(t, err)

		_, err = tarWriter.Write([]byte(file.content))
		require.NoError(t, err)
	}

	err := tarWriter.Close()
	require.NoError(t, err)

	// Compress with xz
	var xzBuf bytes.Buffer

	xzWriter, err := xz.NewWriter(&xzBuf)
	require.NoError(t, err)

	_, err = xzWriter.Write(tarBuf.Bytes())
	require.NoError(t, err)
	err = xzWriter.Close()
	require.NoError(t, err)

	// Test listing contents
	result, err := ListTarxzContents(&xzBuf)
	require.NoError(t, err)

	expected := []string{"hyperfine/hyperfine", "hyperfine/README.md"}
	assert.Equal(t, expected, result)
}

func TestListArchiveContents(t *testing.T) {
	tests := []struct {
		name      string
		assetName string
		setupFunc func() *bytes.Buffer
		expected  []string
	}{
		{
			name:      "tar.gz file",
			assetName: "hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.gz",
			setupFunc: func() *bytes.Buffer {
				var buf bytes.Buffer
				gzWriter := gzip.NewWriter(&buf)
				tarWriter := tar.NewWriter(gzWriter)

				header := &tar.Header{
					Name: "hyperfine",
					Mode: 0o755,
					Size: 100,
				}
				tarWriter.WriteHeader(header)
				tarWriter.Write([]byte(strings.Repeat("x", 100)))
				tarWriter.Close()
				gzWriter.Close()

				return &buf
			},
			expected: []string{"hyperfine"},
		},
		{
			name:      "zip file",
			assetName: "hyperfine-v1.16.1-x86_64-pc-windows-msvc.zip",
			setupFunc: func() *bytes.Buffer {
				var buf bytes.Buffer
				zipWriter := zip.NewWriter(&buf)
				writer, _ := zipWriter.Create("hyperfine.exe")
				writer.Write([]byte("binary"))
				zipWriter.Close()

				return &buf
			},
			expected: []string{"hyperfine.exe"},
		},
		{
			name:      "unsupported format",
			assetName: "hyperfine.deb",
			setupFunc: func() *bytes.Buffer {
				return &bytes.Buffer{}
			},
			expected: nil,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			buf := testCase.setupFunc()
			result, err := ListArchiveContents(testCase.assetName, buf)

			if testCase.expected == nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported archive format")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, result)
			}
		})
	}
}
