package web

import (
	"bytes"
	"github.com/kpkym/koe/cmd/web/config"
	"github.com/kpkym/koe/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
)

func initConfig[T any](key string) *T {
	t := new(T)
	buffer := bytes.NewBufferString(configString)
	v := viper.New()
	v.SetConfigType("toml")
	v.ReadConfig(buffer)

	var err error
	if key != "" {
		err = v.UnmarshalKey(key, t)
	} else {
		err = v.Unmarshal(t)
	}

	if err != nil {
		os.Exit(1)
	}

	return t
}

func initPlag(flagConfig *config.FlagConfig) {
	flags := Cmd.Flags()
	flags.StringVarP(&flagConfig.Port, "port", "p", utils.IgnoreErr(homedir.Expand(flagConfig.Port)), "服务器端口地址")
}
