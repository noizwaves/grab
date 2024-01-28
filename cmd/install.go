package cmd

import (
	"fmt"

	"github.com/noizwaves/grab/pkg"
	"github.com/spf13/cobra"
)

func makeInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:          "install",
		Short:        "Install missing dependencies",
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

			err = pkg.Install(context)
			if err != nil {
				return fmt.Errorf("error installing: %w", err)
			}

			return nil
		},
	}

	return installCmd
}
