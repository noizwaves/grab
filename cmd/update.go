package cmd //nolint:dupl

import (
	"fmt"
	"os"

	"github.com/noizwaves/grab/pkg"
	"github.com/noizwaves/grab/pkg/github"
	"github.com/spf13/cobra"
)

func makeUpdateCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:          "update [package-name]",
		Short:        "Updates packages to use latest remote version",
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

			updater := pkg.Updater{
				GitHubClient: github.NewClient(),
			}

			var packageName string
			if len(args) > 0 {
				packageName = args[0]
			}

			err = updater.Update(gCtx, packageName, os.Stdout)
			if err != nil {
				return fmt.Errorf("error upgrading: %w", err)
			}

			return nil
		},
	}

	return updateCmd
}
