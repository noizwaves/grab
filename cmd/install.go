package cmd

import (
	"fmt"

	"github.com/noizwaves/dotlocalbin/pkg"
	"github.com/spf13/cobra"
)

func makeInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install missing dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := pkg.NewContext()
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
