package cmd

import (
	"fmt"

	"github.com/noizwaves/garb/pkg"
	"github.com/spf13/cobra"
)

func makeUpdateCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Updates packages to use latest remote version",
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := pkg.NewContext()
			if err != nil {
				return fmt.Errorf("error loading context: %w", err)
			}

			err = pkg.Update(context)
			if err != nil {
				return fmt.Errorf("error upgrading: %w", err)
			}

			return nil
		},
	}

	return updateCmd
}
