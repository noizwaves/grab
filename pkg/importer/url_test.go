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
			name: "basic GitHub URL",
			url:  "https://github.com/boyter/scc",
			want: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				Original:     "https://github.com/boyter/scc",
			},
			wantErr: false,
		},
		{
			name: "GitHub URL with path",
			url:  "https://github.com/boyter/scc/releases/tag/v3.5.0",
			want: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				Original:     "https://github.com/boyter/scc/releases/tag/v3.5.0",
			},
			wantErr: false,
		},
		{
			name: "GitHub URL with issues path",
			url:  "https://github.com/boyter/scc/issues",
			want: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				Original:     "https://github.com/boyter/scc/issues",
			},
			wantErr: false,
		},
		{
			name: "GitHub URL with trailing slash",
			url:  "https://github.com/example/repo/",
			want: &GitHubReleaseURL{
				Organization: "example",
				Repository:   "repo",
				Original:     "https://github.com/example/repo/",
			},
			wantErr: false,
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
			url:     "http://github.com/boyter/scc",
			wantErr: true,
		},
		{
			name:    "non-GitHub domain",
			url:     "https://gitlab.com/boyter/scc",
			wantErr: true,
		},
		{
			name:    "GitHub subdomain not allowed",
			url:     "https://api.github.com/boyter/scc",
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
			name: "basic URL",
			url: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
			},
			want: "boyter/scc",
		},
		{
			name: "another URL",
			url: &GitHubReleaseURL{
				Organization: "example",
				Repository:   "repo",
			},
			want: "example/repo",
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
			name: "valid URL",
			releaseURL: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
			},
			wantErr: false,
		},
		{
			name: "missing organization",
			releaseURL: &GitHubReleaseURL{
				Repository: "scc",
			},
			wantErr: true,
		},
		{
			name: "missing repository",
			releaseURL: &GitHubReleaseURL{
				Organization: "boyter",
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
