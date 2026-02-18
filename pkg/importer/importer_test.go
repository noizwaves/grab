package importer

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"os"
	"path"
	"testing"

	"github.com/noizwaves/grab/pkg"
	"github.com/noizwaves/grab/pkg/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock GitHub client for testing.
type MockGitHubClient struct {
	latestRelease     *github.Release
	latestReleaseErr  error
	downloadResponses map[string][]byte
	downloadErrors    map[string]error
}

func (m *MockGitHubClient) GetLatestRelease(_, _ string) (*github.Release, error) {
	if m.latestRelease != nil {
		return m.latestRelease, nil
	}

	if m.latestReleaseErr != nil {
		return nil, m.latestReleaseErr
	}

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
		"linux,amd64": "hyperfine-v{{ .Version }}-x86_64-unknown-linux-gnu.tar.gz",
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

	result, err := detectEmbeddedBinaryPaths(
		mockClient, "sharkdp", "hyperfine", release, "hyperfine", detectedAssets, "1.16.1",
	)
	require.NoError(t, err)
	require.NotNil(t, result)

	expected := map[string]string{
		"linux,amd64": "hyperfine-v{{ .Version }}-x86_64-unknown-linux-gnu/hyperfine",
	}

	assert.Equal(t, expected, *result)
}

func TestDetectEmbeddedBinaryPathsWithSubdirectory(t *testing.T) {
	release := &github.Release{
		TagName: "v1.16.1",
	}

	detectedAssets := map[string]string{
		"linux,amd64": "hyperfine-v{{ .Version }}-x86_64-unknown-linux-gnu.tar.gz",
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

	result, err := detectEmbeddedBinaryPaths(
		mockClient, "sharkdp", "hyperfine", release, "hyperfine", detectedAssets, "1.16.1",
	)
	require.NoError(t, err)
	require.NotNil(t, result)

	expected := map[string]string{
		"linux,amd64": "hyperfine-v{{ .Version }}-x86_64-unknown-linux-gnu/hyperfine",
	}

	assert.Equal(t, expected, *result)
}

func TestDetectEmbeddedBinaryPathsSkipsNonArchives(t *testing.T) {
	release := &github.Release{
		TagName: "v1.16.1",
	}

	detectedAssets := map[string]string{
		"linux,amd64": "hyperfine-v{{ .Version }}-x86_64-unknown-linux-gnu", // Not an archive
	}

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{},
		downloadErrors:    map[string]error{},
	}

	result, err := detectEmbeddedBinaryPaths(
		mockClient, "sharkdp", "hyperfine", release, "hyperfine", detectedAssets, "1.16.1",
	)
	require.NoError(t, err)

	// Should return nil since non-archive assets are skipped
	assert.Nil(t, result)
}

func TestDetectEmbeddedBinaryPathsHandlesDownloadFailure(t *testing.T) {
	release := &github.Release{
		TagName: "v1.16.1",
	}

	detectedAssets := map[string]string{
		"linux,amd64": "hyperfine-v{{ .Version }}-x86_64-unknown-linux-gnu.tar.gz",
	}

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{},
		downloadErrors: map[string]error{
			"hyperfine-v1.16.1-x86_64-unknown-linux-gnu.tar.gz": errors.New("download failed"),
		},
	}

	result, err := detectEmbeddedBinaryPaths(
		mockClient, "sharkdp", "hyperfine", release, "hyperfine", detectedAssets, "1.16.1",
	)

	// Should return error when download fails
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to download asset")
}

func TestDetectEmbeddedBinaryPathsVersionTemplating(t *testing.T) {
	release := &github.Release{
		TagName: "v2.5.0",
	}

	detectedAssets := map[string]string{
		"darwin,amd64": "tool-v{{ .Version }}-darwin-amd64.tar.gz",
		"linux,arm64":  "tool-v{{ .Version }}-linux-arm64.tar.gz",
	}

	// Create archives with version in binary path
	darwinTarGz := createTestTarGz(map[string]string{
		"tool-v2.5.0-darwin-amd64/bin/tool": "darwin binary content",
		"tool-v2.5.0-darwin-amd64/README":   "readme",
	})

	linuxTarGz := createTestTarGz(map[string]string{
		"tool-v2.5.0-linux-arm64/bin/tool": "linux binary content",
		"tool-v2.5.0-linux-arm64/LICENSE":  "license",
	})

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{
			"tool-v2.5.0-darwin-amd64.tar.gz": darwinTarGz,
			"tool-v2.5.0-linux-arm64.tar.gz":  linuxTarGz,
		},
		downloadErrors: map[string]error{},
	}

	result, err := detectEmbeddedBinaryPaths(mockClient, "example", "tool", release, "tool", detectedAssets, "2.5.0")
	require.NoError(t, err)
	require.NotNil(t, result)

	expected := map[string]string{
		"darwin,amd64": "tool-v{{ .Version }}-darwin-amd64/bin/tool",
		"linux,arm64":  "tool-v{{ .Version }}-linux-arm64/bin/tool",
	}

	assert.Equal(t, expected, *result)
}

func TestDetectEmbeddedBinaryPathsVersionTemplatingWithoutVersionInPath(t *testing.T) {
	release := &github.Release{
		TagName: "v1.0.0",
	}

	detectedAssets := map[string]string{
		"linux,amd64": "simple-tool-linux.tar.gz",
	}

	// Create archive without version in binary path - should not be templated
	linuxTarGz := createTestTarGz(map[string]string{
		"bin/simple-tool": "binary content",
		"docs/README":     "readme",
	})

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{
			"simple-tool-linux.tar.gz": linuxTarGz,
		},
		downloadErrors: map[string]error{},
	}

	result, err := detectEmbeddedBinaryPaths(
		mockClient, "example", "simple-tool", release, "simple-tool", detectedAssets, "1.0.0",
	)
	require.NoError(t, err)
	require.NotNil(t, result)

	expected := map[string]string{
		"linux,amd64": "bin/simple-tool", // No templating since no version in path
	}

	assert.Equal(t, expected, *result)
}

func TestDetectEmbeddedBinaryPathsVersionDetectionError(t *testing.T) {
	release := &github.Release{
		TagName: "invalid-tag", // This should cause version detection to fail
	}

	detectedAssets := map[string]string{
		"linux,amd64": "tool-1.0.0-linux.tar.gz",
	}

	// Create archive with version in binary path
	linuxTarGz := createTestTarGz(map[string]string{
		"tool-1.0.0-linux/tool": "binary content",
	})

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{
			"tool-1.0.0-linux.tar.gz": linuxTarGz,
		},
		downloadErrors: map[string]error{},
	}

	result, err := detectEmbeddedBinaryPaths(mockClient, "example", "tool", release, "tool", detectedAssets, "")
	require.NoError(t, err)
	require.NotNil(t, result)

	// Should return literal path without templating when version detection fails
	expected := map[string]string{
		"linux,amd64": "tool-1.0.0-linux/tool",
	}

	assert.Equal(t, expected, *result)
}

func TestDetectEmbeddedBinaryPathsWithCustomPackageName(t *testing.T) {
	release := &github.Release{
		TagName: "v1.0.0",
	}

	detectedAssets := map[string]string{
		"linux,amd64": "custom-tool-v{{ .Version }}-linux.tar.gz",
	}

	// Create archive with binary matching custom package name, not repo name
	linuxTarGz := createTestTarGz(map[string]string{
		"custom-tool-v1.0.0-linux/my-custom-name": "binary content",
		"custom-tool-v1.0.0-linux/README.md":      "readme",
	})

	mockClient := &MockGitHubClient{
		downloadResponses: map[string][]byte{
			"custom-tool-v1.0.0-linux.tar.gz": linuxTarGz,
		},
		downloadErrors: map[string]error{},
	}

	// Use custom package name instead of default repo name
	customPackageName := "my-custom-name"
	result, err := detectEmbeddedBinaryPaths(
		mockClient, "example", "repo-name", release, customPackageName, detectedAssets, "1.0.0",
	)
	require.NoError(t, err)
	require.NotNil(t, result)

	expected := map[string]string{
		"linux,amd64": "custom-tool-v{{ .Version }}-linux/my-custom-name",
	}

	assert.Equal(t, expected, *result)
}

func makeEmptyGrabContext(t *testing.T) *pkg.GrabContext {
	t.Helper()

	dir := t.TempDir()

	err := os.MkdirAll(path.Join(dir, ".grab", "repository"), 0o755)
	if err != nil {
		t.Fatal(err)
	}

	emptyConfig := []byte(`packages: {}`)

	err = os.WriteFile(path.Join(dir, ".grab", "config.yml"), emptyConfig, 0o644)
	if err != nil {
		t.Fatal(err)
	}

	gCtx, err := pkg.NewGrabContext(path.Join(dir, ".grab"), path.Join(dir, "bin"))
	if err != nil {
		t.Fatal(err)
	}

	return gCtx
}

func TestImportPackageReturnsResult(t *testing.T) {
	gCtx := makeEmptyGrabContext(t)

	// Create tar.gz archives for each platform
	linuxAmd64TarGz := createTestTarGz(map[string]string{
		"mytool": "binary content",
	})
	linuxArm64TarGz := createTestTarGz(map[string]string{
		"mytool": "binary content",
	})
	darwinAmd64TarGz := createTestTarGz(map[string]string{
		"mytool": "binary content",
	})
	darwinArm64TarGz := createTestTarGz(map[string]string{
		"mytool": "binary content",
	})

	release := &github.Release{
		TagName: "v1.2.3",
		Assets: []github.Asset{
			{Name: "mytool-1.2.3-linux-amd64.tar.gz"},
			{Name: "mytool-1.2.3-linux-arm64.tar.gz"},
			{Name: "mytool-1.2.3-darwin-amd64.tar.gz"},
			{Name: "mytool-1.2.3-darwin-arm64.tar.gz"},
		},
	}

	mockClient := &MockGitHubClient{
		latestRelease: release,
		downloadResponses: map[string][]byte{
			"mytool-1.2.3-linux-amd64.tar.gz":  linuxAmd64TarGz,
			"mytool-1.2.3-linux-arm64.tar.gz":  linuxArm64TarGz,
			"mytool-1.2.3-darwin-amd64.tar.gz": darwinAmd64TarGz,
			"mytool-1.2.3-darwin-arm64.tar.gz": darwinArm64TarGz,
		},
		downloadErrors: map[string]error{},
	}

	imp := NewImporter(mockClient)

	result, err := imp.ImportPackage(gCtx, "https://github.com/example/mytool", "", &bytes.Buffer{})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "mytool", result.PackageName)
	assert.Equal(t, "1.2.3", result.Version)
}

func TestImportPackageWithCustomName(t *testing.T) {
	gCtx := makeEmptyGrabContext(t)

	linuxAmd64TarGz := createTestTarGz(map[string]string{
		"my-custom-tool": "binary content",
	})
	linuxArm64TarGz := createTestTarGz(map[string]string{
		"my-custom-tool": "binary content",
	})
	darwinAmd64TarGz := createTestTarGz(map[string]string{
		"my-custom-tool": "binary content",
	})
	darwinArm64TarGz := createTestTarGz(map[string]string{
		"my-custom-tool": "binary content",
	})

	release := &github.Release{
		TagName: "v2.0.0",
		Assets: []github.Asset{
			{Name: "mytool-2.0.0-linux-amd64.tar.gz"},
			{Name: "mytool-2.0.0-linux-arm64.tar.gz"},
			{Name: "mytool-2.0.0-darwin-amd64.tar.gz"},
			{Name: "mytool-2.0.0-darwin-arm64.tar.gz"},
		},
	}

	mockClient := &MockGitHubClient{
		latestRelease: release,
		downloadResponses: map[string][]byte{
			"mytool-2.0.0-linux-amd64.tar.gz":  linuxAmd64TarGz,
			"mytool-2.0.0-linux-arm64.tar.gz":  linuxArm64TarGz,
			"mytool-2.0.0-darwin-amd64.tar.gz": darwinAmd64TarGz,
			"mytool-2.0.0-darwin-arm64.tar.gz": darwinArm64TarGz,
		},
		downloadErrors: map[string]error{},
	}

	imp := NewImporter(mockClient)

	result, err := imp.ImportPackage(gCtx, "https://github.com/example/mytool", "my-custom-tool", &bytes.Buffer{})
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "my-custom-tool", result.PackageName)
	assert.Equal(t, "2.0.0", result.Version)
}
