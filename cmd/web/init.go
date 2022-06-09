package web

import (
	"bytes"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kpkym/koe/cmd/web/config"
	"github.com/kpkym/koe/colly"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
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
	logrus.Info("初始化目录树")
	tree := BuildTree()
	cache.NewMapCache[[]*others.Node]().Set("tree", tree)

	// cleanUp(tree)
}

func cleanUp(tree []*others.Node) {
	utils.GoSafe(func() {
		imageDir := filepath.Join(global.GetServiceContext().Config.FlagConfig.DataDir, "imgs")
		table := global.GetServiceContext().DB.Table("work_domains").Session(&gorm.Session{})

		// 检查db图片完整性 图片小于3个需要删除db数据重新抓取
		keys := make(map[string]int)
		var lt3Codes []string
		for _, entry := range utils.IgnoreErr(ioutil.ReadDir(imageDir)) {
			if cs := utils.ListCode(entry.Name()); len(cs) != 0 {
				keys[cs[0]]++
			}
		}
		for k, v := range keys {
			if v != 3 {
				lt3Codes = append(lt3Codes, k)
			}
		}
		table.Where("code IN ?", lt3Codes).Delete(&domain.WorkDomain{})

		// 不在目录树的 需要删除db数据和图片
		var dbCodes []string
		table.Select("code").Scan(&dbCodes)
		fileTreeCodesSet := hashset.New(utils.Map[string, any](utils.ListMyCode(tree), utils.Str2Any)...)
		dbCodesSet := hashset.New(utils.Map[string, any](dbCodes, utils.Str2Any)...)

		needCrawlCodesSet := fileTreeCodesSet.Difference(dbCodesSet).Values()
		if needRemoveCodesSet := dbCodesSet.Difference(fileTreeCodesSet).Values(); len(needRemoveCodesSet) > 0 {
			// 需要删除的
			table.Where("code IN ?", needRemoveCodesSet).Delete(&domain.WorkDomain{})
			for _, f := range utils.IgnoreErr(ioutil.ReadDir(imageDir)) {
				for _, rmCode := range needRemoveCodesSet {
					if strings.Contains(f.Name(), fmt.Sprint(rmCode)) {
						os.Remove(filepath.Join(imageDir, f.Name()))
					}
				}
			}
		}

		if needCrawlCodes := utils.Map[any, string](needCrawlCodesSet, utils.Any2Str); len(needCrawlCodes) > 0 {
			logrus.Info("爬虫抓取", needCrawlCodes)
			colly.C(needCrawlCodes, func(workDomain *domain.WorkDomain) {
				global.GetServiceContext().DB.Create(workDomain)
			})
		}
		logrus.Info("爬虫完成")
	}, "爬虫出错: ", string(debug.Stack()))
}
