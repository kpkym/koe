package web

import (
	"bytes"
	"github.com/jinzhu/copier"
	"github.com/kpkym/koe/cmd/web/config"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/model/pb"
	"github.com/kpkym/koe/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func initConfig[T any](key string, t *T) {
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
}

func initPlag(flagConfig *config.FlagConfig) {
	flags := Cmd.Flags()
	flags.StringVarP(&flagConfig.Port, "port", "p", utils.IgnoreErr(homedir.Expand(flagConfig.Port)), "服务器端口地址")
	flags.StringVar(&flagConfig.Proxy, "proxy", flagConfig.Proxy, "代理地址: http://127.0.0.1:abcd")
	// flags.StringVar(&flagConfig.NasCacheFile, "nasFile", utils.IgnoreErr(homedir.Expand(flagConfig.NasCacheFile)), "NAS缓存文件")
	flags.StringVar(&flagConfig.DataDir, "dataDir", utils.GetFileBaseOnPwd(flagConfig.DataDir), "服务器数据文件夹")
	flags.StringVar(&flagConfig.ScanDir, "koeDir", utils.IgnoreErr(homedir.Expand(flagConfig.ScanDir)), "扫描文件夹")
}

func InitTree() {
	cacheHolder := pb.PBNode{
		Children: make([]*pb.PBNode, 0),
	}
	copier.Copy(&cacheHolder.Children, utils.BuildTree())
	logrus.Info("初始化目录树")
	cache.Set[*pb.PBNode]("tree", &cacheHolder)
}
