package cmd

import (
	"fmt"
	"os"

	"github.com/noizwaves/grab/pkg"
	"github.com/noizwaves/grab/pkg/github"
	"github.com/noizwaves/grab/pkg/importer"
	"github.com/spf13/cobra"
)

func makeGetCommand() *cobra.Command {
	var packageName string

	getCmd := &cobra.Command{
		Use:   "get [GITHUB_REPO_URL]",
		Short: "Import, configure, and install a package in one step",
		Long: `
Automatically import a package from a GitHub repo URL, add it to the config, and install the binary.

This combines the import, config edit, and install steps into a single command.

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
			return runGetCommand(packageName, args)
		},
	}

	getCmd.Flags().StringVarP(
		&packageName, "name", "n", "",
		"Override package name (must be lowercase with no whitespace)",
	)

	return getCmd
}

func runGetCommand(packageName string, args []string) error {
	if packageName != "" {
		err := validatePackageName(packageName)
		if err != nil {
			return fmt.Errorf("invalid package name: %w", err)
		}
	}

	inputURL := args[0]

	err := validateGitHubRepoURL(inputURL)
	if err != nil {
		return fmt.Errorf("invalid GitHub release URL: %w", err)
	}

	gCtx, err := newGrabContext()
	if err != nil {
		return fmt.Errorf("error loading context: %w", err)
	}

	imp := importer.NewImporter(github.NewClient())

	result, err := imp.ImportPackage(gCtx, inputURL, packageName, os.Stdout)
	if err != nil {
		return fmt.Errorf("error importing: %w", err)
	}

	err = gCtx.AddPackageToConfig(result.PackageName, result.Version)
	if err != nil {
		return fmt.Errorf("error adding package to config: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Added %s@%s to config\n", result.PackageName, result.Version)

	gCtx, err = newGrabContext()
	if err != nil {
		return fmt.Errorf("error reloading context: %w", err)
	}

	installer := pkg.Installer{
		GitHubClient: github.NewClient(),
	}

	err = installer.Install(gCtx, result.PackageName, os.Stdout)
	if err != nil {
		return fmt.Errorf("error installing: %w", err)
	}

	return nil
}
