package web

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/cmd/web/config"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/router"
	"github.com/kpkym/koe/utils"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

var (
	Cmd = &cobra.Command{
		Use:   "web",
		Short: "启动web服务",
		Run: func(_ *cobra.Command, _ []string) {
			os.MkdirAll(global.GetServiceContext().Config.DataDir, 0700)
			global.AddDB(initDB())
			InitTree()
			web()
		},
	}
)

//go:embed config/config.toml
var configString string

//go:embed dist
var dist embed.FS

func init() {
	flagConfig := &config.FlagConfig{}
	commonConfig := &config.CommonConfig{}
	initConfig("flag", flagConfig)
	initConfig("common", commonConfig)

	global.SetServiceContext(&config.Config{
		FlagConfig:   flagConfig,
		CommonConfig: commonConfig,
	})

	initPlag(flagConfig)
}

func web() {
	serve := router.GetGinServe()

	fs := http.FS(utils.IgnoreErr(fs.Sub(dist, "dist")))

	serve.NoRoute(func(context *gin.Context) {
		context.FileFromFS("/"+context.Request.RequestURI, fs)
	})

	serve.Run(":" + global.GetServiceContext().Config.FlagConfig.Port)
}

func initDB() *gorm.DB {
	// 初始化数据库
	db := utils.IgnoreErr(gorm.Open(sqlite.Open(filepath.Join(global.GetServiceContext().Config.FlagConfig.DataDir, "koe.db")),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)}))
	// 迁移 schema
	db.AutoMigrate(&domain.WorkDomain{})
	return db
}
