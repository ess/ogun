package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string
	Build   string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Report the ogun version",
	Long: `Report the ogun version

The version and build time are set during our release build, and this can
help you to guarantee that you're using the most recent version.`,

	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func showVersion() {
	fmt.Printf("ogun %s (Build %s)\n", Version, Build)
}
