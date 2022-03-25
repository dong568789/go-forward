package cmd

import (
	"github.com/dong568789/go-forward/library/conf"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "tcp forward",
	Long:  `a tcp forward tool`,
}

var cfgFile string

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./conf.yml", "config file (default is project root dir /conf.yml)")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	conf.Init(cfgFile)
}
