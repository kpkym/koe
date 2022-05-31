package web

import (
	"bytes"
	_ "embed"
	"fmt"
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
	"path/filepath"
)

var (
	c   = initConfig()
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
	// 初始化配置
	c = initConfig()

	Cmd.Flags().StringVarP(&global.NasCacheFile, "nasFile", "n", utils.IgnoreErr(homedir.Expand(c.Common.NasCacheFile)), "NAS缓存文件")
	Cmd.Flags().StringVarP(&global.SqliteDataFile, "db", "d", utils.IgnoreErr(homedir.Expand(c.Common.SqliteDataFile)), "sqlite数据库文件")
	Cmd.Flags().StringVarP(&global.ScanDir, "koe", "k", utils.IgnoreErr(homedir.Expand(c.Common.ScanDir)), "扫描文件夹")
}

func web() {
	global.SetServiceContext(initDB(), c)

	serve := router.GetGinServe()
	serve.StaticFS("/static", http.Dir(global.ScanDir))
	serve.Group("/file").GET("/cover/:type/:id", func(c *gin.Context) {
		imgPath := filepath.Join(utils.GetFileBaseOnPwd("data, imgs"),
			filepath.Base(utils.GetImgUrl(c.Param("id"), c.Param("type"))))

		logrus.Infof("查找图片: %s", imgPath)
		c.File(imgPath)
	})

	serve.Run(":" + c.Common.Port)
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

func initConfig() *config.Config {
	buffer := bytes.NewBufferString(configString)

	v := viper.New()
	v.SetConfigType("toml")

	v.ReadConfig(buffer)

	conf := &config.Config{}
	err := v.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
	return conf
}
