package cmd

import (
	"fmt"
	"os"

	"github.com/noizwaves/grab/pkg"
	"github.com/noizwaves/grab/pkg/github"
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

			installer := pkg.Installer{
				GitHubClient: github.NewClient(),
			}

			err = installer.Install(context, os.Stdout)
			if err != nil {
				return fmt.Errorf("error installing: %w", err)
			}

			return nil
		},
	}

	return installCmd
}
