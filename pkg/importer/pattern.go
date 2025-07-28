package importer

import (
	"errors"
	"regexp"
)

type DetectedPattern struct {
	Value string
}

type PatternDetector struct {
	patterns []pattern
}

type pattern struct {
	Name  string
	Regex *regexp.Regexp
	Value string
}

func (pd *PatternDetector) AnalyzeOne(value string) (*DetectedPattern, error) {
	for _, pattern := range pd.patterns {
		if pattern.Regex.MatchString(value) {
			return &DetectedPattern{Value: pattern.Value}, nil
		}
	}

	return nil, errors.New("no matching pattern found")
}

// Detecting release names
func NewReleaseNamePatternDetector() *PatternDetector {
	patterns := []pattern{
		{
			Name:  "v-prefixed",
			Regex: regexp.MustCompile(`^v\d+\.\d+\.\d+$`),
			Value: "v{{ .Version }}",
		},
		{
			Name:  "version-only",
			Regex: regexp.MustCompile(`^\d+\.\d+\.\d+$`),
			Value: "{{ .Version }}",
		},
	}

	return &PatternDetector{
		patterns: patterns,
	}
}

// Detecting platform
var darwinPlatformPatterns = []pattern{
	{
		Name:  "darwin",
		Regex: regexp.MustCompile(`(?i)darwin`),
		Value: "darwin",
	},
}

var linuxPlatformPatterns = []pattern{
	{
		Name:  "linux",
		Regex: regexp.MustCompile(`(?i)linux`),
		Value: "linux",
	},
}

func NewPlatformPatternDetector(platform string) *PatternDetector {
	if platform == "darwin" {
		return &PatternDetector{
			patterns: darwinPlatformPatterns,
		}
	} else if platform == "linux" {
		return &PatternDetector{
			patterns: linuxPlatformPatterns,
		}
	}

	panic("unsupported platform: " + platform)
}

// Detecting architecture
var amd64ArchitecturePatterns = []pattern{
	{
		Name:  "amd64",
		Regex: regexp.MustCompile(`(?i)amd64`),
		Value: "amd64",
	},
	{
		Name:  "x86_64",
		Regex: regexp.MustCompile(`(?i)x86_64`),
		Value: "x86_64",
	},
}

var arm64ArchitecturePatterns = []pattern{
	{
		Name:  "arm64",
		Regex: regexp.MustCompile(`(?i)arm64`),
		Value: "arm64",
	},
	{
		Name:  "aarch64",
		Regex: regexp.MustCompile(`(?i)aarch64`),
		Value: "aarch64",
	},
}

func NewArchitecturePatternDetector(architecture string) *PatternDetector {
	if architecture == "amd64" {
		return &PatternDetector{
			patterns: amd64ArchitecturePatterns,
		}
	} else if architecture == "arm64" {
		return &PatternDetector{
			patterns: arm64ArchitecturePatterns,
		}
	}

	panic("unsupported architecture: " + architecture)
}
