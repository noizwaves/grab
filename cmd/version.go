package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v0.0.0"

func makeVersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:          "version",
		Short:        "Print grab version",
		SilenceUsage: true,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(version)
		},
	}

	return versionCmd
}
