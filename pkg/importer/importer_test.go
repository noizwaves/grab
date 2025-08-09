package importer

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"testing"

	"github.com/noizwaves/grab/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock GitHub client for testing.
type MockGitHubClient struct {
	downloadResponses map[string][]byte
	downloadErrors    map[string]error
}

func (m *MockGitHubClient) GetLatestRelease(_, _ string) (*github.Release, error) {
	return nil, errors.New("not implemented for test")
}

func (m *MockGitHubClient) GetReleaseByTag(_, _, _ string) (*github.Release, error) {
	return nil, errors.New("not implemented for test")
}

func (m *MockGitHubClient) DownloadReleaseAsset(_, _, _, asset string) ([]byte, error) {
	if err, exists := m.downloadErrors[asset]; exists {
		return nil, err
	}

	if data, exists := m.downloadResponses[asset]; exists {
		return data, nil
	}

	return nil, errors.New("asset not found in mock")
}

func TestIsArchiveAsset(t *testing.T) {
	tests := []struct {
		assetName string
		expected  bool
	}{
		{"hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.gz", true},
		{"hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tgz", true},
		{"hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.xz", true},
		{"hyperfine-v1.16.1-x86_64-pc-windows-msvc.zip", true},
		{"hyperfine-v1.16.1-x86_64-unknown-linux-gnu", false},
		{"hyperfine_1.16.1_amd64.deb", false},
		{"hyperfine-1.16.1-1.x86_64.rpm", false},
	}

	for _, tt := range tests {
		t.Run(tt.assetName, func(t *testing.T) {
			result := isArchiveAsset(tt.assetName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindBinaryInArchive(t *testing.T) {
	tests := []struct {
		name        string
		files       []string
		packageName string
		expected    string
		expectError bool
	}{
		{
			name:        "exact match at root",
			files:       []string{"hyperfine", "README.md", "LICENSE"},
			packageName: "hyperfine",
			expected:    "hyperfine",
			expectError: false,
		},
		{
			name:        "exact match in subdirectory",
			files:       []string{"hyperfine/hyperfine", "hyperfine/README.md", "hyperfine/LICENSE"},
			packageName: "hyperfine",
			expected:    "hyperfine/hyperfine",
			expectError: false,
		},
		{
			name:        "exact match with extension on Windows",
			files:       []string{"hyperfine.exe", "README.md"},
			packageName: "hyperfine.exe",
			expected:    "hyperfine.exe",
			expectError: false,
		},
		{
			name:        "partial match fallback",
			files:       []string{"hyperfine-bin", "README.md"},
			packageName: "hyperfine",
			expected:    "hyperfine-bin",
			expectError: false,
		},
		{
			name:        "no match",
			files:       []string{"other-binary", "README.md"},
			packageName: "hyperfine",
			expected:    "",
			expectError: true,
		},
		{
			name:        "skip directories",
			files:       []string{"hyperfine/", "hyperfine/hyperfine", "README.md"},
			packageName: "hyperfine",
			expected:    "hyperfine/hyperfine",
			expectError: false,
		},
		{
			name:        "prefer exact over partial match",
			files:       []string{"hyperfine-extended", "hyperfine", "README.md"},
			packageName: "hyperfine",
			expected:    "hyperfine",
			expectError: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := findBinaryInArchive(testCase.files, testCase.packageName)

			if testCase.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, result)
			}
		})
	}
}

func createTestTarGz(files map[string]string) []byte {
	var buf bytes.Buffer

	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	for name, content := range files {
		header := &tar.Header{
			Name: name,
			Mode: 0o755,
			Size: int64(len(content)),
		}
		tarWriter.WriteHeader(header)
		tarWriter.Write([]byte(content))
	}

	tarWriter.Close()
	gzWriter.Close()

	return buf.Bytes()
}

func TestDetectEmbeddedBinaryPaths(t *testing.T) {
	release := &github.Release{
		TagName: "v1.16.1",
	}

	detectedAssets := map[string]string{
		"linux,amd64": "hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.gz",
	}

	// Create tar.gz with binary in subdirectory (not at root, so it needs embedded path)
	linuxTarGz := createTestTarGz(map[string]string{
		"hyperfine-v1.16.1-x86_64-unknown-linux-gnu/hyperfine": "binary content",
		"hyperfine-v1.16.1-x86_64-unknown-linux-gnu/README.md": "readme",
	})

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{
			"hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.gz": linuxTarGz,
		},
		downloadErrors: map[string]error{},
	}

	result, err := detectEmbeddedBinaryPaths(mockClient, "sharkdp", "hyperfine", release, "hyperfine", detectedAssets)
	require.NoError(t, err)
	require.NotNil(t, result)

	expected := map[string]string{
		"linux,amd64": "hyperfine-v1.16.1-x86_64-unknown-linux-gnu/hyperfine",
	}

	assert.Equal(t, expected, *result)
}

func TestDetectEmbeddedBinaryPathsWithSubdirectory(t *testing.T) {
	release := &github.Release{
		TagName: "v1.16.1",
	}

	detectedAssets := map[string]string{
		"linux,amd64": "hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.gz",
	}

	// Create tar.gz with binary in subdirectory (like real hyperfine releases)
	linuxTarGz := createTestTarGz(map[string]string{
		"hyperfine-v1.16.1-x86_64-unknown-linux-gnu/hyperfine": "binary content",
		"hyperfine-v1.16.1-x86_64-unknown-linux-gnu/README.md": "readme",
		"hyperfine-v1.16.1-x86_64-unknown-linux-gnu/LICENSE":   "license",
	})

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{
			"hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.gz": linuxTarGz,
		},
		downloadErrors: map[string]error{},
	}

	result, err := detectEmbeddedBinaryPaths(mockClient, "sharkdp", "hyperfine", release, "hyperfine", detectedAssets)
	require.NoError(t, err)
	require.NotNil(t, result)

	expected := map[string]string{
		"linux,amd64": "hyperfine-v1.16.1-x86_64-unknown-linux-gnu/hyperfine",
	}

	assert.Equal(t, expected, *result)
}

func TestDetectEmbeddedBinaryPathsSkipsNonArchives(t *testing.T) {
	release := &github.Release{
		TagName: "v1.16.1",
	}

	detectedAssets := map[string]string{
		"linux,amd64": "hyperfine-v1.16.1-x86_64-unknown-linux-gnu", // Not an archive
	}

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{},
		downloadErrors:    map[string]error{},
	}

	result, err := detectEmbeddedBinaryPaths(mockClient, "sharkdp", "hyperfine", release, "hyperfine", detectedAssets)
	require.NoError(t, err)

	// Should return nil since non-archive assets are skipped
	assert.Nil(t, result)
}

func TestDetectEmbeddedBinaryPathsHandlesDownloadFailure(t *testing.T) {
	release := &github.Release{
		TagName: "v1.16.1",
	}

	detectedAssets := map[string]string{
		"linux,amd64": "hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.gz",
	}

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{},
		downloadErrors: map[string]error{
			"hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.gz": errors.New("download failed"),
		},
	}

	result, err := detectEmbeddedBinaryPaths(mockClient, "sharkdp", "hyperfine", release, "hyperfine", detectedAssets)

	// Should return error when download fails
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to download asset")
}
