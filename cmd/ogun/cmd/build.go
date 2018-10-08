package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ess/ogun/pkg/ogun/fs"

	"github.com/ess/ogun/cmd/ogun/workflows"
)

// 2018 08 29 22 43 47

var buildCmd = &cobra.Command{
	Use:   "build <application name>",
	Short: "Build a portable application package",
	Long: `Build a portable application package

Given an application name, build the application and generate an all-inclusive
tarball for distribution.`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("no application given")
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(ReleaseName) == 0 {
			ReleaseName = workflows.GenerateBuildNumber()
		}

		logfile, err := fs.CreateBuildLog(args[0], ReleaseName)
		if err != nil {
			return fmt.Errorf("could not open log")
		}

		defer logfile.Close()

		Logger.AddOutput(logfile)

		return workflows.Perform(
			&workflows.BuildingAnApp{
				ApplicationName: args[0],
				ReleaseName:     ReleaseName,
				Apps:            fs.NewApplicationService(Logger),
				Packs:           fs.NewBuildpackService(Logger),
				Releases:        fs.NewReleaseService(Logger),
				Logger:          Logger,
			},
		)
	},

	SilenceUsage:  true,
	SilenceErrors: true,
}

var ReleaseName string

func init() {
	buildCmd.Flags().StringVarP(
		&ReleaseName,
		"release",
		"r",
		"",
		"the name of the release to build",
	)

	RootCmd.AddCommand(buildCmd)
}
