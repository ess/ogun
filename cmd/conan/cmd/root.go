package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "conan",
	Short: "Conan the Deployer",
	Long:  `Conan the Deployer`,
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
	viper.AddConfigPath("/etc/conan")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
	}
}
