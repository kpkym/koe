package web

import (
	"bytes"
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kpk-koe/cmd/web/config"
	"kpk-koe/global"
	"kpk-koe/model/domain"
	"kpk-koe/router"
	"kpk-koe/utils"
	"net/http"
	"os"
	"path/filepath"
)

var (
	Cmd = &cobra.Command{
		Use:   "web",
		Short: "启动web服务",
		Run: func(_ *cobra.Command, _ []string) {
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

	global.SetServiceContext(initDB(), &config.Config{
		FlagConfig:   *flagConfig,
		CommonConfig: *commonConfig,
	})

	flags := Cmd.Flags()
	flags.StringVarP(&flagConfig.Port, "port", "p", utils.IgnoreErr(homedir.Expand(flagConfig.Port)), "NAS缓存文件")
	flags.StringVarP(&flagConfig.NasCacheFile, "nasFile", "n", utils.IgnoreErr(homedir.Expand(flagConfig.NasCacheFile)), "NAS缓存文件")
	flags.StringVarP(&flagConfig.SqliteDataFile, "db", "d", utils.IgnoreErr(homedir.Expand(flagConfig.SqliteDataFile)), "sqlite数据库文件")
	flags.StringVarP(&flagConfig.ScanDir, "koe", "k", utils.IgnoreErr(homedir.Expand(flagConfig.ScanDir)), "扫描文件夹")
}

func web() {
	serve := router.GetGinServe()
	serve.StaticFS("/static", http.Dir(global.ScanDir))
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
	db := utils.IgnoreErr(gorm.Open(sqlite.Open(global.SqliteDataFile), &gorm.Config{}))
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
