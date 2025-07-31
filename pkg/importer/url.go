package importer

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// GitHubReleaseURL represents a parsed GitHub URL.
type GitHubReleaseURL struct {
	// Organization is the GitHub organization or user name
	Organization string

	// Repository is the GitHub repository name
	Repository string

	// Original is the original URL string
	Original string
}

// String returns a string representation of the parsed URL.
func (u *GitHubReleaseURL) String() string {
	return fmt.Sprintf("%s/%s", u.Organization, u.Repository)
}

// Validate checks if the parsed URL components are valid.
func (u *GitHubReleaseURL) Validate() error {
	if u.Organization == "" {
		return errors.New("organization cannot be empty")
	}

	if u.Repository == "" {
		return errors.New("repository cannot be empty")
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
	minURLPathComponents = 2
)

// ParseGitHubReleaseURL parses a GitHub URL and extracts org/repo components.
func ParseGitHubReleaseURL(inputURL string) (*GitHubReleaseURL, error) {
	parsedURL, err := parseAndValidateURL(inputURL)
	if err != nil {
		return nil, err
	}

	organization, repository, err := extractOrgAndRepo(inputURL, parsedURL.Path)
	if err != nil {
		return nil, err
	}

	result := &GitHubReleaseURL{
		Organization: organization,
		Repository:   repository,
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

func extractOrgAndRepo(inputURL, path string) (string, string, error) {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < minURLPathComponents {
		return "", "", newParseError(inputURL, "path", "insufficient path components for GitHub URL")
	}

	return pathParts[0], pathParts[1], nil
}
