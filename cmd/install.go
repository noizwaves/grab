package cmd //nolint:dupl

import (
	"fmt"
	"os"

	"github.com/noizwaves/grab/pkg"
	"github.com/noizwaves/grab/pkg/github"
	"github.com/spf13/cobra"
)

func makeInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:          "install [package-name]",
		Short:        "Install missing dependencies",
		Args:         cobra.MaximumNArgs(1),
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

			installer := pkg.Installer{
				GitHubClient: github.NewClient(),
			}

			var packageName string
			if len(args) > 0 {
				packageName = args[0]
			}

			err = installer.Install(gCtx, packageName, os.Stdout)
			if err != nil {
				return fmt.Errorf("error installing: %w", err)
			}

			return nil
		},
	}

	return installCmd
}
