package importer

import (
	"testing"

	"github.com/noizwaves/grab/pkg/github"
)

func TestNewPatternDetector(t *testing.T) {
	detector := NewPatternDetector()
	if detector == nil {
		t.Fatal("Expected detector to be created, got nil")
	}

	if detector.versionRegex == nil {
		t.Error("Expected version regex to be initialized")
	}

	if len(detector.patterns) == 0 {
		t.Error("Expected patterns to be initialized")
	}
}

func TestPatternType_String(t *testing.T) {
	tests := []struct {
		name     string
		pType    PatternType
		expected string
	}{
		{"version-only", PatternTypeVersionOnly, "version-only"},
		{"v-prefix", PatternTypeVPrefix, "v-prefix"},
		{"release-prefix", PatternTypeReleasePrefix, "release-prefix"},
		{"custom", PatternTypeCustom, "custom"},
		{"unknown", PatternType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pType.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPatternDetector_AnalyzeReleaseNames_VPrefix(t *testing.T) {
	detector := NewPatternDetector()

	releaseNames := []string{
		"v1.0.0",
		"v1.1.0",
		"v1.2.0",
		"v2.0.0",
	}

	pattern, err := detector.AnalyzeReleaseNames(releaseNames)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pattern.Type != PatternTypeVPrefix {
		t.Errorf("Expected pattern type %v, got %v", PatternTypeVPrefix, pattern.Type)
	}

	if pattern.Template != "v{{ .Version }}" {
		t.Errorf("Expected template 'v{{ .Version }}', got '%s'", pattern.Template)
	}

	if pattern.Confidence != 1.0 {
		t.Errorf("Expected confidence 1.0, got %f", pattern.Confidence)
	}

	if pattern.Matches != 4 {
		t.Errorf("Expected 4 matches, got %d", pattern.Matches)
	}

	if pattern.Total != 4 {
		t.Errorf("Expected total 4, got %d", pattern.Total)
	}
}

func TestPatternDetector_AnalyzeReleaseNames_VersionOnly(t *testing.T) {
	detector := NewPatternDetector()

	releaseNames := []string{
		"1.0.0",
		"1.1.0",
		"1.2.0",
		"2.0.0",
	}

	pattern, err := detector.AnalyzeReleaseNames(releaseNames)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pattern.Type != PatternTypeVersionOnly {
		t.Errorf("Expected pattern type %v, got %v", PatternTypeVersionOnly, pattern.Type)
	}

	if pattern.Template != "{{ .Version }}" {
		t.Errorf("Expected template '{{ .Version }}', got '%s'", pattern.Template)
	}

	if pattern.Confidence != 1.0 {
		t.Errorf("Expected confidence 1.0, got %f", pattern.Confidence)
	}
}

func TestPatternDetector_AnalyzeReleaseNames_ReleasePrefix(t *testing.T) {
	detector := NewPatternDetector()

	releaseNames := []string{
		"release-1.0.0",
		"release-1.1.0",
		"release-1.2.0",
		"release-2.0.0",
	}

	pattern, err := detector.AnalyzeReleaseNames(releaseNames)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pattern.Type != PatternTypeReleasePrefix {
		t.Errorf("Expected pattern type %v, got %v", PatternTypeReleasePrefix, pattern.Type)
	}

	if pattern.Template != "release-{{ .Version }}" {
		t.Errorf("Expected template 'release-{{ .Version }}', got '%s'", pattern.Template)
	}

	if pattern.Confidence != 1.0 {
		t.Errorf("Expected confidence 1.0, got %f", pattern.Confidence)
	}
}

func TestPatternDetector_AnalyzeReleaseNames_CustomPattern(t *testing.T) {
	detector := NewPatternDetector()

	releaseNames := []string{
		"myapp-1.0.0-final",
		"myapp-1.1.0-final",
		"myapp-1.2.0-final",
		"myapp-2.0.0-final",
	}

	pattern, err := detector.AnalyzeReleaseNames(releaseNames)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pattern.Type != PatternTypeCustom {
		t.Errorf("Expected pattern type %v, got %v", PatternTypeCustom, pattern.Type)
	}

	if pattern.Template != "myapp-{{ .Version }}-final" {
		t.Errorf("Expected template 'myapp-{{ .Version }}-final', got '%s'", pattern.Template)
	}

	if pattern.Confidence != 1.0 {
		t.Errorf("Expected confidence 1.0, got %f", pattern.Confidence)
	}

	if pattern.Prefix != "myapp-" {
		t.Errorf("Expected prefix 'myapp-', got '%s'", pattern.Prefix)
	}

	if pattern.Suffix != "-final" {
		t.Errorf("Expected suffix '-final', got '%s'", pattern.Suffix)
	}
}

func TestPatternDetector_AnalyzeReleaseNames_MixedPatterns(t *testing.T) {
	detector := NewPatternDetector()

	releaseNames := []string{
		"v1.0.0",
		"v1.1.0",
		"1.2.0", // Different pattern
		"v2.0.0",
	}

	pattern, err := detector.AnalyzeReleaseNames(releaseNames)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should detect v-prefix as the dominant pattern (3/4 matches)
	if pattern.Type != PatternTypeVPrefix {
		t.Errorf("Expected pattern type %v, got %v", PatternTypeVPrefix, pattern.Type)
	}

	if pattern.Confidence != 0.75 { // 3/4 = 0.75
		t.Errorf("Expected confidence 0.75, got %f", pattern.Confidence)
	}

	if pattern.Matches != 3 {
		t.Errorf("Expected 3 matches, got %d", pattern.Matches)
	}
}

func TestPatternDetector_AnalyzeReleaseNames_PreReleaseVersions(t *testing.T) {
	detector := NewPatternDetector()

	releaseNames := []string{
		"v1.0.0-alpha",
		"v1.0.0-beta",
		"v1.0.0",
		"v1.1.0-alpha.1",
	}

	pattern, err := detector.AnalyzeReleaseNames(releaseNames)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pattern.Type != PatternTypeVPrefix {
		t.Errorf("Expected pattern type %v, got %v", PatternTypeVPrefix, pattern.Type)
	}

	if pattern.Confidence != 1.0 {
		t.Errorf("Expected confidence 1.0, got %f", pattern.Confidence)
	}
}

func TestPatternDetector_AnalyzeReleaseNames_EmptyInput(t *testing.T) {
	detector := NewPatternDetector()

	_, err := detector.AnalyzeReleaseNames([]string{})
	if err == nil {
		t.Fatal("Expected error for empty input, got nil")
	}

	if err.Error() != "no release names to analyze" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestPatternDetector_AnalyzeReleaseNames_NoPattern(t *testing.T) {
	detector := NewPatternDetector()

	releaseNames := []string{
		"invalid-name-1",
		"another-invalid-2",
		"completely-different",
		"no-version-here",
	}

	_, err := detector.AnalyzeReleaseNames(releaseNames)
	if err == nil {
		t.Fatal("Expected error for no detectable pattern, got nil")
	}

	if err.Error() != "no consistent pattern detected" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestPatternDetector_AnalyzeReleases(t *testing.T) {
	detector := NewPatternDetector()

	releases := []github.Release{
		{Name: "v1.0.0"},
		{Name: "v1.1.0"},
		{Name: "v1.2.0"},
		{Name: "v2.0.0"},
	}

	pattern, err := detector.AnalyzeReleases(releases)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pattern.Type != PatternTypeVPrefix {
		t.Errorf("Expected pattern type %v, got %v", PatternTypeVPrefix, pattern.Type)
	}
}

func TestPatternDetector_AnalyzeReleases_EmptyInput(t *testing.T) {
	detector := NewPatternDetector()

	_, err := detector.AnalyzeReleases([]github.Release{})
	if err == nil {
		t.Fatal("Expected error for empty input, got nil")
	}

	if err.Error() != "no releases to analyze" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestPatternDetector_ValidatePattern(t *testing.T) {
	detector := NewPatternDetector()

	tests := []struct {
		name        string
		pattern     *DetectedPattern
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil pattern",
			pattern:     nil,
			expectError: true,
			errorMsg:    "pattern cannot be nil",
		},
		{
			name: "low confidence",
			pattern: &DetectedPattern{
				Confidence: 0.5,
				Value:      "v{{ .Version }}",
			},
			expectError: true,
			errorMsg:    "pattern confidence too low",
		},
		{
			name: "missing version placeholder",
			pattern: &DetectedPattern{
				Confidence: 0.9,
				Value:      "v1.2.3",
			},
			expectError: true,
			errorMsg:    "template must contain {{ .Version }} placeholder",
		},
		{
			name: "valid pattern",
			pattern: &DetectedPattern{
				Confidence: 0.9,
				Value:      "v{{ .Version }}",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := detector.ValidatePattern(tt.pattern, []string{})

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) &&
		(s == substr || s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
