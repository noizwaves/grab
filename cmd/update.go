package cmd

import (
	"fmt"
	"os"

	"github.com/noizwaves/grab/pkg"
	"github.com/noizwaves/grab/pkg/github"
	"github.com/spf13/cobra"
)

func makeUpdateCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:          "update",
		Short:        "Updates packages to use latest remote version",
		SilenceUsage: true,
		PreRun: func(_ *cobra.Command, _ []string) {
			err := configureLogging()
			cobra.CheckErr(err)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			context, err := newContext()
			if err != nil {
				return fmt.Errorf("error loading context: %w", err)
			}

			updater := pkg.Updater{
				GitHubClient: github.NewClient(),
			}

			err = updater.Update(context, os.Stdout)
			if err != nil {
				return fmt.Errorf("error upgrading: %w", err)
			}

			return nil
		},
	}

	return updateCmd
}
