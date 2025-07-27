package importer

import (
	"fmt"
	"net/url"
	"strings"
)

// GitHubReleaseURL represents a parsed GitHub release URL
type GitHubReleaseURL struct {
	// Organization is the GitHub organization or user name
	Organization string
	
	// Repository is the GitHub repository name
	Repository string
	
	// Tag is the release tag (e.g., "v1.2.3", "1.2.3")
	// Empty for latest releases
	Tag string
	
	// IsLatest indicates if this is a /releases/latest URL
	IsLatest bool
	
	// Original is the original URL string
	Original string
}

// String returns a string representation of the parsed URL
func (u *GitHubReleaseURL) String() string {
	if u.IsLatest {
		return fmt.Sprintf("%s/%s (latest)", u.Organization, u.Repository)
	}
	return fmt.Sprintf("%s/%s@%s", u.Organization, u.Repository, u.Tag)
}

// Validate checks if the parsed URL components are valid
func (u *GitHubReleaseURL) Validate() error {
	if u.Organization == "" {
		return fmt.Errorf("organization cannot be empty")
	}
	if u.Repository == "" {
		return fmt.Errorf("repository cannot be empty")
	}
	if !u.IsLatest && u.Tag == "" {
		return fmt.Errorf("tag cannot be empty for non-latest releases")
	}
	return nil
}

// ParseError represents an error that occurred during URL parsing
type ParseError struct {
	URL    string
	Reason string
	Field  string // Which field/component failed validation
}

func (e *ParseError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("invalid %s in URL %q: %s", e.Field, e.URL, e.Reason)
	}
	return fmt.Sprintf("invalid URL %q: %s", e.URL, e.Reason)
}

// newParseError creates a new ParseError
func newParseError(url, field, reason string) *ParseError {
	return &ParseError{
		URL:    url,
		Field:  field,
		Reason: reason,
	}
}

// ParseGitHubReleaseURL parses a GitHub release URL and extracts components
func ParseGitHubReleaseURL(inputURL string) (*GitHubReleaseURL, error) {
	// Parse the URL using net/url
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, newParseError(inputURL, "format", fmt.Sprintf("invalid URL format: %v", err))
	}

	// Validate scheme
	if parsedURL.Scheme != "https" {
		return nil, newParseError(inputURL, "scheme", fmt.Sprintf("URL must use HTTPS scheme, got: %s", parsedURL.Scheme))
	}

	// Validate host
	if parsedURL.Host != "github.com" {
		return nil, newParseError(inputURL, "host", fmt.Sprintf("URL must be from github.com, got: %s", parsedURL.Host))
	}

	// Extract path components
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 4 {
		return nil, newParseError(inputURL, "path", "insufficient path components for GitHub release URL")
	}

	// Extract organization and repository (first two components)
	organization := pathParts[0]
	repository := pathParts[1]

	// Validate that this is a releases URL
	if pathParts[2] != "releases" {
		return nil, newParseError(inputURL, "path", fmt.Sprintf("expected 'releases', got '%s'", pathParts[2]))
	}

	// Parse release type and tag
	var tag string
	var isLatest bool

	releaseType := pathParts[3]
	switch releaseType {
	case "latest":
		isLatest = true
	case "tag":
		if len(pathParts) < 5 {
			return nil, newParseError(inputURL, "tag", "tag release URL must specify a version")
		}
		tag = pathParts[4]
	default:
		// Handle direct version format: /releases/version
		if len(pathParts) == 4 {
			tag = pathParts[3]
		} else {
			return nil, newParseError(inputURL, "path", fmt.Sprintf("expected 'tag', 'latest', or direct version, got '%s'", releaseType))
		}
	}

	// Create the parsed URL struct
	result := &GitHubReleaseURL{
		Organization: organization,
		Repository:   repository,
		Tag:          tag,
		IsLatest:     isLatest,
		Original:     inputURL,
	}

	// Validate the parsed components
	if err := result.Validate(); err != nil {
		return nil, newParseError(inputURL, "components", err.Error())
	}

	return result, nil
}