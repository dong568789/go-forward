package conf

import (
	"github.com/dong568789/go-forward/library/util"
	"github.com/spf13/viper"
)

type Configs struct {
	Proxies []string `yaml:"proxies"`
}

var Conf = &Configs{}

func Init(path string) {
	if path == "" || !util.Exists(path) {
		util.Log().Panic("config file not exists: ", viper.ConfigFileUsed())
	}
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		util.Log().Info("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(Conf)
	if err != nil {
		util.Log().Panic("Unmarshal config fail: %v", err)
	}
}
