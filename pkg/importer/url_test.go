package importer

import (
	"testing"
)

func TestParseGitHubReleaseURL(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGitHubReleaseURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGitHubReleaseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Organization != tt.want.Organization {
					t.Errorf("ParseGitHubReleaseURL() Organization = %v, want %v", got.Organization, tt.want.Organization)
				}
				if got.Repository != tt.want.Repository {
					t.Errorf("ParseGitHubReleaseURL() Repository = %v, want %v", got.Repository, tt.want.Repository)
				}
				if got.Tag != tt.want.Tag {
					t.Errorf("ParseGitHubReleaseURL() Tag = %v, want %v", got.Tag, tt.want.Tag)
				}
				if got.IsLatest != tt.want.IsLatest {
					t.Errorf("ParseGitHubReleaseURL() IsLatest = %v, want %v", got.IsLatest, tt.want.IsLatest)
				}
				if got.Original != tt.want.Original {
					t.Errorf("ParseGitHubReleaseURL() Original = %v, want %v", got.Original, tt.want.Original)
				}
			}
		})
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
		name    string
		url     *GitHubReleaseURL
		wantErr bool
	}{
		{
			name: "valid tag release",
			url: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				Tag:          "v3.5.0",
				IsLatest:     false,
			},
			wantErr: false,
		},
		{
			name: "valid latest release",
			url: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				IsLatest:     true,
			},
			wantErr: false,
		},
		{
			name: "missing organization",
			url: &GitHubReleaseURL{
				Repository: "scc",
				Tag:        "v3.5.0",
			},
			wantErr: true,
		},
		{
			name: "missing repository",
			url: &GitHubReleaseURL{
				Organization: "boyter",
				Tag:          "v3.5.0",
			},
			wantErr: true,
		},
		{
			name: "missing tag for non-latest",
			url: &GitHubReleaseURL{
				Organization: "boyter",
				Repository:   "scc",
				IsLatest:     false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.url.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GitHubReleaseURL.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}