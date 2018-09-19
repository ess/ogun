package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ess/ogun/fs"
	"github.com/ess/ogun/log"

	"github.com/ess/ogun/cmd/ogun/workflows"
)

// 2018 08 29 22 43 47

var buildCmd = &cobra.Command{
	Use:   "build <application>",
	Short: "Build a portable application package",
	Long: `Build a portable application package

Given an application name, build the application and generate an all-inclusive
tarball for distribution.

The name for the specific buildpack to use to build the application can be
provided with the --buildpack flag.

If no specific buildpack is provided, an attempt is made to detect the proper
buildpack based on the application source.`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("no application given")
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		release := workflows.GenerateBuildNumber()
		logger := log.NewLogger()
		logfile, err := fs.CreateBuildLog(args[0], release)
		if err != nil {
			return fmt.Errorf("could not open log")
		}

		defer logfile.Close()

		logger.AddOutput(logfile)

		return workflows.Perform(
			&workflows.BuildingAnApp{
				BuildpackName:   buildPackName,
				ApplicationName: args[0],
				ReleaseName:     release,
				Apps:            fs.NewApplicationService(logger),
				Packs:           fs.NewBuildpackService(logger),
				Releases:        fs.NewReleaseService(logger),
				Logger:          logger,
			},
		)
	},

	SilenceUsage:  true,
	SilenceErrors: true,
}

var buildPackName string

func init() {
	buildCmd.Flags().StringVar(&buildPackName, "buildpack", "detect",
		"The buildpack to build the application")

	RootCmd.AddCommand(buildCmd)
}
