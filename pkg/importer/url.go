package importer

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// GitHubReleaseURL represents a parsed GitHub release URL.
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

// String returns a string representation of the parsed URL.
func (u *GitHubReleaseURL) String() string {
	if u.IsLatest {
		return fmt.Sprintf("%s/%s (latest)", u.Organization, u.Repository)
	}

	return fmt.Sprintf("%s/%s@%s", u.Organization, u.Repository, u.Tag)
}

// Validate checks if the parsed URL components are valid.
func (u *GitHubReleaseURL) Validate() error {
	if u.Organization == "" {
		return errors.New("organization cannot be empty")
	}

	if u.Repository == "" {
		return errors.New("repository cannot be empty")
	}

	if !u.IsLatest && u.Tag == "" {
		return errors.New("tag cannot be empty for non-latest releases")
	}

	return nil
}

// ParseError represents an error that occurred during URL parsing.
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

// newParseError creates a new ParseError.
func newParseError(urlStr, field, reason string) *ParseError {
	return &ParseError{
		URL:    urlStr,
		Field:  field,
		Reason: reason,
	}
}

const (
	minURLPathComponents = 4
	minTagPathComponents = 5
	pathOrgIndex         = 0
	pathRepoIndex        = 1
	pathReleasesIndex    = 2
	pathTypeIndex        = 3
	pathTagIndex         = 4
)

// ParseGitHubReleaseURL parses a GitHub release URL and extracts components.
func ParseGitHubReleaseURL(inputURL string) (*GitHubReleaseURL, error) {
	parsedURL, err := parseAndValidateURL(inputURL)
	if err != nil {
		return nil, err
	}

	pathParts, err := extractAndValidatePathParts(inputURL, parsedURL.Path)
	if err != nil {
		return nil, err
	}

	tag, isLatest, err := parseReleaseInfo(inputURL, pathParts)
	if err != nil {
		return nil, err
	}

	result := &GitHubReleaseURL{
		Organization: pathParts[pathOrgIndex],
		Repository:   pathParts[pathRepoIndex],
		Tag:          tag,
		IsLatest:     isLatest,
		Original:     inputURL,
	}

	err = result.Validate()
	if err != nil {
		return nil, newParseError(inputURL, "components", err.Error())
	}

	return result, nil
}

func parseAndValidateURL(inputURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, newParseError(inputURL, "format", fmt.Sprintf("invalid URL format: %v", err))
	}

	if parsedURL.Scheme != "https" {
		return nil, newParseError(inputURL, "scheme", "URL must use HTTPS scheme, got: "+parsedURL.Scheme)
	}

	if parsedURL.Host != "github.com" {
		return nil, newParseError(inputURL, "host", "URL must be from github.com, got: "+parsedURL.Host)
	}

	return parsedURL, nil
}

func extractAndValidatePathParts(inputURL, path string) ([]string, error) {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < minURLPathComponents {
		return nil, newParseError(inputURL, "path", "insufficient path components for GitHub release URL")
	}

	if pathParts[pathReleasesIndex] != "releases" {
		return nil, newParseError(inputURL, "path",
			fmt.Sprintf("expected 'releases', got '%s'", pathParts[pathReleasesIndex]))
	}

	return pathParts, nil
}

func parseReleaseInfo(inputURL string, pathParts []string) (string, bool, error) {
	releaseType := pathParts[pathTypeIndex]

	switch releaseType {
	case "latest":
		return "", true, nil
	case "tag":
		if len(pathParts) < minTagPathComponents {
			return "", false, newParseError(inputURL, "tag", "tag release URL must specify a version")
		}

		return pathParts[pathTagIndex], false, nil
	default:
		// Handle direct version format: /releases/version
		if len(pathParts) == minURLPathComponents {
			return pathParts[pathTypeIndex], false, nil
		}

		return "", false, newParseError(inputURL, "path",
			"expected 'tag', 'latest', or direct version, got '"+releaseType+"'")
	}
}
