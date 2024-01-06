package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func makeRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dotlocalbin",
		Short: "User centric dotfile dependency manager",
	}

	rootCmd.AddCommand(makeInstallCommand())

	return rootCmd
}

func Execute() {
	rootCmd := makeRootCommand()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
