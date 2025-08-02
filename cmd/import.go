package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/noizwaves/grab/pkg/github"
	"github.com/noizwaves/grab/pkg/importer"
	"github.com/spf13/cobra"
)

func makeImportCommand() *cobra.Command {
	importCmd := &cobra.Command{
		Use:          "import <github-release-url>",
		Short:        "Generate package spec from GitHub release",
		Long:         "Automatically generate a package specification YAML file by analyzing a GitHub release URL",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		PreRun: func(_ *cobra.Command, _ []string) {
			err := configureLogging()
			cobra.CheckErr(err)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			gCtx, err := newGrabContext()
			if err != nil {
				return fmt.Errorf("error loading context: %w", err)
			}

			inputURL := args[0]

			// Validate GitHub release URL
			err = validateGitHubReleaseURL(inputURL)
			if err != nil {
				return fmt.Errorf("invalid GitHub release URL: %w", err)
			}

			importer := importer.NewImporter(github.NewClient())

			err = importer.Import(gCtx, inputURL, os.Stdout)
			if err != nil {
				return fmt.Errorf("error installing: %w", err)
			}

			return nil
		},
	}

	return importCmd
}

const (
	minPathComponents    = 4
	minTagPathComponents = 5
)

// validateGitHubReleaseURL validates that the URL is a valid GitHub release URL.
// Valid URL schemes:
// - https://github.com/<org>/<repo>/releases/tag/<version>
// - https://github.com/<org>/<repo>/releases/latest
// - https://github.com/<org>/<repo>/releases/<version>
func validateGitHubReleaseURL(inputURL string) error {
	parsedURL, err := validateBasicURL(inputURL)
	if err != nil {
		return err
	}

	return validateGitHubReleasePath(parsedURL.Path)
}

func validateBasicURL(inputURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("URL must use HTTPS scheme, got: %s", parsedURL.Scheme)
	}

	if parsedURL.Host != "github.com" {
		return nil, fmt.Errorf("URL must be from github.com, got: %s", parsedURL.Host)
	}

	return parsedURL, nil
}

func validateGitHubReleasePath(path string) error {
	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < minPathComponents {
		return errors.New("invalid GitHub path structure")
	}

	if pathParts[2] != "releases" {
		return errors.New("URL must be a GitHub releases URL")
	}

	return validateReleaseFormat(pathParts)
}

func validateReleaseFormat(pathParts []string) error {
	switch pathParts[3] {
	case "latest":
		return nil
	case "tag":
		if len(pathParts) < minTagPathComponents {
			return errors.New("tag URL must specify a version")
		}

		return nil
	default:
		if len(pathParts) == minPathComponents {
			return nil
		}

		return errors.New("URL must be a release tag, latest release, or direct version")
	}
}
