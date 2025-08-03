package importer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectVersionRegex(t *testing.T) {
	tests := []struct {
		name     string
		tagName  string
		expected string
	}{
		{
			name:     "returns default version regex for any tag name",
			tagName:  "v1.2.3",
			expected: `\d+\.\d+\.\d+`,
		},
		{
			name:     "returns default version regex for version-only tag",
			tagName:  "1.2.3",
			expected: `\d+\.\d+\.\d+`,
		},
		{
			name:     "returns default version regex for complex tag",
			tagName:  "release-v1.2.3-beta",
			expected: `\d+\.\d+\.\d+`,
		},
		{
			name:     "returns default version regex for empty tag",
			tagName:  "",
			expected: `\d+\.\d+\.\d+`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := DetectVersionRegex(test.tagName)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestDetectVersionValue(t *testing.T) {
	tests := []struct {
		name         string
		versionRegex string
		tagName      string
		expected     string
		expectError  bool
	}{
		{
			name:         "extracts version from v-prefixed tag",
			versionRegex: `\d+\.\d+\.\d+`,
			tagName:      "v1.2.3",
			expected:     "1.2.3",
			expectError:  false,
		},
		{
			name:         "extracts version from version-only tag",
			versionRegex: `\d+\.\d+\.\d+`,
			tagName:      "1.2.3",
			expected:     "1.2.3",
			expectError:  false,
		},
		{
			name:         "extracts version from complex tag",
			versionRegex: `\d+\.\d+\.\d+`,
			tagName:      "release-v1.2.3-beta",
			expected:     "1.2.3",
			expectError:  false,
		},
		{
			name:         "extracts version with patch version zero",
			versionRegex: `\d+\.\d+\.\d+`,
			tagName:      "v2.1.0",
			expected:     "2.1.0",
			expectError:  false,
		},
		{
			name:         "extracts first version when multiple present",
			versionRegex: `\d+\.\d+\.\d+`,
			tagName:      "v1.2.3-to-v2.0.0",
			expected:     "1.2.3",
			expectError:  false,
		},
		{
			name:         "returns error when no version found",
			versionRegex: `\d+\.\d+\.\d+`,
			tagName:      "no-version-here",
			expected:     "",
			expectError:  true,
		},
		{
			name:         "returns error for empty tag name",
			versionRegex: `\d+\.\d+\.\d+`,
			tagName:      "",
			expected:     "",
			expectError:  true,
		},
		{
			name:         "works with custom version regex",
			versionRegex: `v\d+\.\d+`,
			tagName:      "release-v1.2-stable",
			expected:     "v1.2",
			expectError:  false,
		},
		{
			name:         "handles semantic version with build metadata",
			versionRegex: `\d+\.\d+\.\d+`,
			tagName:      "v1.2.3+build.123",
			expected:     "1.2.3",
			expectError:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := DetectVersionValue(test.versionRegex, test.tagName)

			if test.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestUnrenderVersionValue(t *testing.T) {
	tests := []struct {
		name           string
		value          string
		versionLiteral string
		expected       string
	}{
		{
			name:           "replaces version literal in filename",
			value:          "crush_1.2.3.tar.gz",
			versionLiteral: "1.2.3",
			expected:       "crush_{{ .Version }}.tar.gz",
		},
		{
			name:           "replaces multiple occurrences of version literal",
			value:          "app-1.2.3-1.2.3.tar.gz",
			versionLiteral: "1.2.3",
			expected:       "app-{{ .Version }}-{{ .Version }}.tar.gz",
		},
		{
			name:           "returns original value when version literal not found",
			value:          "crush.tar.gz",
			versionLiteral: "1.2.3",
			expected:       "crush.tar.gz",
		},
		{
			name:           "handles version literal at start of value",
			value:          "1.2.3-release.tar.gz",
			versionLiteral: "1.2.3",
			expected:       "{{ .Version }}-release.tar.gz",
		},
		{
			name:           "handles version literal at end of value",
			value:          "release-1.2.3",
			versionLiteral: "1.2.3",
			expected:       "release-{{ .Version }}",
		},
		{
			name:           "handles version literal as entire value",
			value:          "1.2.3",
			versionLiteral: "1.2.3",
			expected:       "{{ .Version }}",
		},
		{
			name:           "handles empty version literal",
			value:          "crush_1.2.3.tar.gz",
			versionLiteral: "",
			expected:       "crush_1.2.3.tar.gz",
		},
		{
			name:           "handles empty value",
			value:          "",
			versionLiteral: "1.2.3",
			expected:       "",
		},
		{
			name:           "handles different version formats",
			value:          "app-v2.1.0.zip",
			versionLiteral: "v2.1.0",
			expected:       "app-{{ .Version }}.zip",
		},
		{
			name:           "handles version with patch zero",
			value:          "binary-2.1.0-linux.tar.gz",
			versionLiteral: "2.1.0",
			expected:       "binary-{{ .Version }}-linux.tar.gz",
		},
		{
			name:           "handles version in path-like structure",
			value:          "releases/1.2.3/app.tar.gz",
			versionLiteral: "1.2.3",
			expected:       "releases/{{ .Version }}/app.tar.gz",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := UnrenderVersionValue(test.value, test.versionLiteral)
			assert.Equal(t, test.expected, result)
		})
	}
}
