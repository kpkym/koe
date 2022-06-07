package web

import (
	"bytes"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/jinzhu/copier"
	"github.com/kpkym/koe/cmd/web/config"
	"github.com/kpkym/koe/colly"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/pb"
	"github.com/kpkym/koe/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"runtime/debug"
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
	logrus.Info("初始化目录树")
	tree := utils.BuildTree()

	copier.Copy(&cacheHolder.Children, tree)
	cache.Set[*pb.PBNode]("tree", &cacheHolder)

	go func() {
		defer func() {
			if p := recover(); p != nil {
				logrus.Error("爬虫出错: ", string(debug.Stack()))
			}
		}()

		var dbCodes []string
		global.GetServiceContext().DB.Table("work_domains").Select("code").Scan(&dbCodes)

		codes := utils.Map[string, any](utils.ListMyCode(tree), utils.Str2Any)

		needCrawlCodesSet := hashset.New(codes...).Difference(hashset.New(utils.Map[string, any](dbCodes, utils.Str2Any)...)).Values()
		if needCrawlCodes := utils.Map[any, string](needCrawlCodesSet, utils.Any2Str); len(needCrawlCodes) > 0 {
			logrus.Info("爬虫抓取", needCrawlCodes)
			c, _ := colly.C(needCrawlCodes)

			v := make([]*domain.WorkDomain, 0, len(c))
			for _, value := range c {
				v = append(v, value)
			}
			global.GetServiceContext().DB.Create(&v)
		}
	}()

}
