package cmd

import (
	"fmt"
	"os"

	"github.com/noizwaves/grab/pkg"
	"github.com/spf13/cobra"
)

func makeUpdateCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:          "update",
		Short:        "Updates packages to use latest remote version",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			err := configureLogging()
			cobra.CheckErr(err)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := newContext()
			if err != nil {
				return fmt.Errorf("error loading context: %w", err)
			}

			err = pkg.Update(context, os.Stdout)
			if err != nil {
				return fmt.Errorf("error upgrading: %w", err)
			}

			return nil
		},
	}

	return updateCmd
}
