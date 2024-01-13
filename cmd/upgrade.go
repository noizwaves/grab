package cmd

import (
	"fmt"

	"github.com/noizwaves/garb/pkg"
	"github.com/spf13/cobra"
)

func makeUpgradeCommand() *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade configured version to latest upstream",
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := pkg.NewContext()
			if err != nil {
				return fmt.Errorf("error loading context: %w", err)
			}

			err = pkg.Upgrade(context)
			if err != nil {
				return fmt.Errorf("error upgrading: %w", err)
			}

			return nil
		},
	}

	return upgradeCmd
}
