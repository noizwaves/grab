package cmd

import (
	"fmt"

	"github.com/noizwaves/dotlocalbin/pkg"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install missing dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		context, err := pkg.NewContext()
		if err != nil {
			return fmt.Errorf("Error loading context: %w", err)
		}

		return pkg.Install(context)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
