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

//nolint:lll
func makeImportCommand() *cobra.Command {
	var packageName string

	importCmd := &cobra.Command{
		Use:   "import [GITHUB_REPO_URL]",
		Short: "Generate package spec from GitHub repo",
		Long: `
Automatically generate a package specification YAML file by analyzing a GitHub repo URL.

Arguments:
  GITHUB_REPO_URL: GitHub repository URL to analyze (e.g., https://github.com/junegunn/fzf)

Flags:
  -n, --name string: Override package name (default: repository name, must be lowercase with no whitespace)
`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		PreRun: func(_ *cobra.Command, _ []string) {
			err := configureLogging()
			cobra.CheckErr(err)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			// Validate package name if provided
			if packageName != "" {
				err := validatePackageName(packageName)
				if err != nil {
					return fmt.Errorf("invalid package name: %w", err)
				}
			}

			gCtx, err := newGrabContext()
			if err != nil {
				return fmt.Errorf("error loading context: %w", err)
			}

			inputURL := args[0]

			// Validate GitHub release URL
			err = validateGitHubRepoURL(inputURL)
			if err != nil {
				return fmt.Errorf("invalid GitHub release URL: %w", err)
			}

			importer := importer.NewImporter(github.NewClient())

			err = importer.Import(gCtx, inputURL, packageName, os.Stdout)
			if err != nil {
				return fmt.Errorf("error installing: %w", err)
			}

			return nil
		},
	}

	importCmd.Flags().StringVarP(&packageName, "name", "n", "", "Override package name (must be lowercase with no whitespace)")

	return importCmd
}

// validateGitHubRepoURL validates that the URL is a valid GitHub URL.
// Valid URL scheme: https://github.com/<org>/<repo>/*
func validateGitHubRepoURL(inputURL string) error {
	parsedURL, err := validateGitHubURL(inputURL)
	if err != nil {
		return err
	}

	const minRequiredPathComponents = 2
	// Check that path has at least owner and repo components
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < minRequiredPathComponents {
		return errors.New("URL must include both owner and repository (e.g., https://github.com/owner/repo)")
	}

	if pathParts[0] == "" || pathParts[1] == "" {
		return errors.New("URL must include both owner and repository (e.g., https://github.com/owner/repo)")
	}

	return nil
}

func validateGitHubURL(inputURL string) (*url.URL, error) {
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

// validatePackageName validates that the package name is lowercase and contains no whitespace.
func validatePackageName(name string) error {
	if strings.TrimSpace(name) != name {
		return errors.New("package name cannot contain leading or trailing whitespace")
	}

	if strings.ContainsAny(name, " \t\n\r") {
		return errors.New("package name cannot contain whitespace characters")
	}

	if name != strings.ToLower(name) {
		return errors.New("package name must be lowercase")
	}

	if name == "" {
		return errors.New("package name cannot be empty")
	}

	return nil
}
