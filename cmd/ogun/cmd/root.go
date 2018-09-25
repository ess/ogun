package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "ogun",
	Short: "Ogun, God of Metalcraft",
	Long: `Ogun, God of Metalcraft
	
Ogun is a utiltiy used, primarily, for crafting portable application releases on
bare metal.

This is the top level of that utility, which does not do very much. Please see
below for subcommands that do much more.`,
}

func Execute() error {
	err := RootCmd.Execute()

	if err != nil {
		fmt.Println(err)
	}

	return err
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/ogun")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
	}
}
