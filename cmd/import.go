package cmd

import (
	"fmt"
	"net/url"
	"strings"

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
			inputURL := args[0]

			// Validate GitHub release URL
			if err := validateGitHubReleaseURL(inputURL); err != nil {
				return fmt.Errorf("invalid GitHub release URL: %w", err)
			}

			// TODO: Implement import logic
			fmt.Printf("Import functionality not yet implemented for URL: %s\n", inputURL)
			return fmt.Errorf("import command is not yet implemented")
		},
	}

	return importCmd
}

// validateGitHubReleaseURL validates that the URL is a valid GitHub release URL
// valid URL schemes:
// - https://github.com/<org>/<repo>/releases/tag/<version>
// - https://github.com/<org>/<repo>/releases/latest
// - https://github.com/<org>/<repo>/releases/<version>
func validateGitHubReleaseURL(inputURL string) error {
	// Parse the URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check scheme
	if parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use HTTPS scheme, got: %s", parsedURL.Scheme)
	}

	// Check host
	if parsedURL.Host != "github.com" {
		return fmt.Errorf("URL must be from github.com, got: %s", parsedURL.Host)
	}

	// Check path structure: /org/repo/releases/tag/version or /org/repo/releases/latest
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 4 {
		return fmt.Errorf("invalid GitHub path structure")
	}

	// Validate path components
	if pathParts[2] != "releases" {
		return fmt.Errorf("URL must be a GitHub releases URL")
	}

	if pathParts[3] == "latest" {
		// /releases/latest format
		return nil
	} else if pathParts[3] == "tag" {
		// /releases/tag/version format
		if len(pathParts) < 5 {
			return fmt.Errorf("tag URL must specify a version")
		}
	} else if len(pathParts) == 4 {
		// /releases/version format - this is valid but redirects to /releases/tag/version
		return nil
	} else {
		return fmt.Errorf("URL must be a release tag, latest release, or direct version")
	}

	return nil
}
