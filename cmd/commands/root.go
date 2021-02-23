package commands

import (
	"github.com/ernie-mlg/ErniePJT-main-api-go/viper"
	"github.com/spf13/cobra"
)

type Command = cobra.Command

func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "ErniePJT",
	Short: "Ernie Project Management Application",
	Long:  `ErniePJT offers workplace messaging across web, PC and phones with archiving, search and integration with your existing systems.`,
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "config.json", "Configuration file to use.")
	RootCmd.PersistentFlags().Bool("disableconfigwatch", false, "When set config.json will not be loaded from disk when the file is changed.")

	viper.BindEnv("config")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
}
