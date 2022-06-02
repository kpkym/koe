package web

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/cmd/web/config"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/router"
	"github.com/kpkym/koe/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
)

var (
	Cmd = &cobra.Command{
		Use:   "web",
		Short: "启动web服务",
		Run: func(_ *cobra.Command, _ []string) {
			global.AddDB(initDB())
			web()
		},
	}
)

//go:embed config/config.toml
var configString string

func init() {
	flagConfig := &config.FlagConfig{}
	commonConfig := &config.CommonConfig{}
	initConfig("flag", flagConfig)
	initConfig("common", commonConfig)

	global.SetServiceContext(&config.Config{
		FlagConfig:   *flagConfig,
		CommonConfig: *commonConfig,
	})

	flags := Cmd.Flags()
	flags.StringVarP(&flagConfig.Port, "port", "p", utils.IgnoreErr(homedir.Expand(flagConfig.Port)), "服务器端口地址")
	flags.StringVar(&flagConfig.Proxy, "proxy", flagConfig.Proxy, "代理地址")

	if flagConfig.Serve != "" {
		flags.StringVarP(&flagConfig.Serve, "serve", "s", flagConfig.Serve, "服务器地址")
	} else {
		flags.StringVarP(&flagConfig.Serve, "serve", "s", fmt.Sprintf("http://%s:%s", utils.GetLocalIp()[0], flagConfig.Port), "服务器地址")
	}

	flags.StringVarP(&flagConfig.NasCacheFile, "nasFile", "n", utils.IgnoreErr(homedir.Expand(flagConfig.NasCacheFile)), "NAS缓存文件")
	flags.StringVar(&flagConfig.SqliteDataFile, "db", utils.IgnoreErr(homedir.Expand(flagConfig.SqliteDataFile)), "sqlite数据库文件")
	flags.StringVar(&flagConfig.ScanDir, "koeDir", utils.IgnoreErr(homedir.Expand(flagConfig.ScanDir)), "扫描文件夹")
}

func web() {
	serve := router.GetGinServe()
	serve.StaticFS("/static", http.Dir(global.GetServiceContext().Config.FlagConfig.ScanDir))
	serve.Group("/file").GET("/cover/:type/:id", func(c *gin.Context) {
		imgPath := filepath.Join(utils.GetFileBaseOnPwd("data, imgs"),
			filepath.Base(utils.GetImgUrl(c.Param("id"), c.Param("type"))))

		logrus.Infof("查找图片: %s", imgPath)
		c.File(imgPath)
	})

	serve.Run(":" + global.GetServiceContext().Config.FlagConfig.Port)
}

func initDB() *gorm.DB {
	// 初始化数据库
	db := utils.IgnoreErr(gorm.Open(sqlite.Open(global.GetServiceContext().Config.FlagConfig.SqliteDataFile), &gorm.Config{}))
	// 迁移 schema
	db.AutoMigrate(&domain.DlDomain{})
	db.AutoMigrate(&domain.TrackDomain{})
	db.AutoMigrate(&domain.FileSystemTreeDomain{})
	return db
}

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
