package cmd

import (
	"github.com/dong568789/go-forward/library"
	"github.com/dong568789/go-forward/library/conf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serverCmd)

	viper.BindPFlags(serverCmd.Flags())
}

var serverCmd = &cobra.Command{
	Use:   "start",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

//执行-次日志文件导入
func run() {
	library.NewProxy().ParseAddr(conf.Conf.Proxies).Run()
}
