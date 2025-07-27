package importer

import (
	"testing"
)

func TestParseGitHubReleaseURL(t *testing.T) {
	tests := getParseTestCases()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ParseGitHubReleaseURL(testCase.url)
			if (err != nil) != testCase.wantErr {
				t.Errorf("ParseGitHubReleaseURL() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !testCase.wantErr {
				validateParseResult(t, got, testCase.want)
			}
		})
	}
}

//nolint:funlen
func getParseTestCases() []struct {
	name    string
	url     string
	want    *GitHubReleaseURL
	wantErr bool
} {
	return []struct {
		name    string
		url     string
		want    *GitHubReleaseURL
		wantErr bool
	}{
		{
			name: "valid tag release URL",
			url:  "https://github.com/boyter/scc/releases/tag/v3.5.0",
			want: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				Tag:          "v3.5.0",
				IsLatest:     false,
				Original:     "https://github.com/boyter/scc/releases/tag/v3.5.0",
			},
			wantErr: false,
		},
		{
			name: "valid latest release URL",
			url:  "https://github.com/boyter/scc/releases/latest",
			want: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				Tag:          "",
				IsLatest:     true,
				Original:     "https://github.com/boyter/scc/releases/latest",
			},
			wantErr: false,
		},
		{
			name: "tag without version prefix",
			url:  "https://github.com/example/repo/releases/tag/1.2.3",
			want: &GitHubReleaseURL{
				Organization: "example",
				Repository:   "repo",
				Tag:          "1.2.3",
				IsLatest:     false,
				Original:     "https://github.com/example/repo/releases/tag/1.2.3",
			},
			wantErr: false,
		},
		{
			name: "direct version format",
			url:  "https://github.com/boyter/scc/releases/v3.5.0",
			want: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				Tag:          "v3.5.0",
				IsLatest:     false,
				Original:     "https://github.com/boyter/scc/releases/v3.5.0",
			},
			wantErr: false,
		},
		{
			name:    "missing tag version",
			url:     "https://github.com/boyter/scc/releases/tag",
			wantErr: true,
		},
		{
			name:    "non-releases URL",
			url:     "https://github.com/boyter/scc/issues",
			wantErr: true,
		},
		{
			name:    "insufficient path components",
			url:     "https://github.com/boyter",
			wantErr: true,
		},
		{
			name:    "invalid URL format",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "HTTP scheme not allowed",
			url:     "http://github.com/boyter/scc/releases/tag/v3.5.0",
			wantErr: true,
		},
		{
			name:    "non-GitHub domain",
			url:     "https://gitlab.com/boyter/scc/releases/tag/v3.5.0",
			wantErr: true,
		},
		{
			name:    "GitHub subdomain not allowed",
			url:     "https://api.github.com/boyter/scc/releases/tag/v3.5.0",
			wantErr: true,
		},
	}
}

func validateParseResult(t *testing.T, got, want *GitHubReleaseURL) {
	t.Helper()

	if got.Organization != want.Organization {
		t.Errorf("ParseGitHubReleaseURL() Organization = %v, want %v", got.Organization, want.Organization)
	}

	if got.Repository != want.Repository {
		t.Errorf("ParseGitHubReleaseURL() Repository = %v, want %v", got.Repository, want.Repository)
	}

	if got.Tag != want.Tag {
		t.Errorf("ParseGitHubReleaseURL() Tag = %v, want %v", got.Tag, want.Tag)
	}

	if got.IsLatest != want.IsLatest {
		t.Errorf("ParseGitHubReleaseURL() IsLatest = %v, want %v", got.IsLatest, want.IsLatest)
	}

	if got.Original != want.Original {
		t.Errorf("ParseGitHubReleaseURL() Original = %v, want %v", got.Original, want.Original)
	}
}

func TestGitHubReleaseURL_String(t *testing.T) {
	tests := []struct {
		name string
		url  *GitHubReleaseURL
		want string
	}{
		{
			name: "tag release",
			url: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				Tag:          "v3.5.0",
				IsLatest:     false,
			},
			want: "boyter/scc@v3.5.0",
		},
		{
			name: "latest release",
			url: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				IsLatest:     true,
			},
			want: "boyter/scc (latest)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.url.String(); got != tt.want {
				t.Errorf("GitHubReleaseURL.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitHubReleaseURL_Validate(t *testing.T) {
	tests := []struct {
		name       string
		releaseURL *GitHubReleaseURL
		wantErr    bool
	}{
		{
			name: "valid tag release",
			releaseURL: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				Tag:          "v3.5.0",
				IsLatest:     false,
			},
			wantErr: false,
		},
		{
			name: "valid latest release",
			releaseURL: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				IsLatest:     true,
			},
			wantErr: false,
		},
		{
			name: "missing organization",
			releaseURL: &GitHubReleaseURL{
				Repository: "scc",
				Tag:        "v3.5.0",
			},
			wantErr: true,
		},
		{
			name: "missing repository",
			releaseURL: &GitHubReleaseURL{
				Organization: "boyter",
				Tag:          "v3.5.0",
			},
			wantErr: true,
		},
		{
			name: "missing tag for non-latest",
			releaseURL: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				IsLatest:     false,
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.releaseURL.Validate()
			if (err != nil) != testCase.wantErr {
				t.Errorf("GitHubReleaseURL.Validate() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}
}
