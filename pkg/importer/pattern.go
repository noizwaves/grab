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

// Detecting release names.
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

// Detecting platform.
func NewPlatformPatternDetector(platform string) *PatternDetector {
	switch platform {
	case "darwin":
		return &PatternDetector{
			patterns: []pattern{
				{
					Name:  "darwin",
					Regex: regexp.MustCompile(`(?i)darwin`),
					Value: "darwin",
				},
			},
		}
	case "linux":
		return &PatternDetector{
			patterns: []pattern{
				{
					Name:  "linux",
					Regex: regexp.MustCompile(`(?i)linux`),
					Value: "linux",
				},
			},
		}
	default:
		panic("unsupported platform: " + platform)
	}
}

// Detecting architecture.
func NewArchitecturePatternDetector(architecture string) *PatternDetector {
	switch architecture {
	case "amd64":
		return &PatternDetector{
			patterns: []pattern{
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
			},
		}
	case "arm64":
		return &PatternDetector{
			patterns: []pattern{
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
			},
		}
	default:
		panic("unsupported architecture: " + architecture)
	}
}

func (pd *PatternDetector) AnalyzeOne(value string) (*DetectedPattern, error) {
	for _, pattern := range pd.patterns {
		if pattern.Regex.MatchString(value) {
			return &DetectedPattern{Value: pattern.Value}, nil
		}
	}

	return nil, errors.New("no matching pattern found")
}
