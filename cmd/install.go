package cmd

import (
	"github.com/noizwaves/dotlocalbin/pkg"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install missing dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkg.Install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
