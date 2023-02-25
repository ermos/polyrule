package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "polyrule",
	Short: "Compile validator rules into multiple languages",
	Long:  `compile validator rules into multiple languages`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	viper.SetDefault("author", "Kilian SMITI <kilian@smiti.fr>")
	viper.SetDefault("license", "MIT")
}

func initConfig() {
	// @TODO get os information ?
}
