package web

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/kpkym/koe/cmd/web/config"
	"github.com/kpkym/koe/colly"
	"github.com/kpkym/koe/dao/cache"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/model/others"
	"github.com/kpkym/koe/router"
	"github.com/kpkym/koe/utils"
	"github.com/kpkym/koe/utils/koe"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/fs"
	"net/http"
	"path/filepath"
)

var (
	Cmd = &cobra.Command{
		Use:   "web",
		Short: "启动web服务",
		Run: func(_ *cobra.Command, _ []string) {
			koe.BuildTree()
			trees, _ := cache.NewMapCache[string, []*others.Node]().Get("trees")
			koe.RemoveIncomplete()
			koe.RemoveNotInTrees(trees, func(needCrawlCodes []string) {
				colly.C(needCrawlCodes, func(workDomain *domain.WorkDomain) {
					global.GetServiceContext().DB.Create(workDomain)
				})
			})
			web()
		},
	}
)

//go:embed config/config.toml
var configString string

//go:embed dist
var dist embed.FS

func init() {
	fcg := initConfig[config.FlagConfig]("flag")
	global.SetServiceContext(&config.Config{
		FlagConfig:   fcg,
		CommonConfig: initConfig[config.CommonConfig]("common"),
	}, initDB())
	initPlag(fcg)
	global.AddSettings(loadSettingsFromDB())
}

func initDB() *gorm.DB {
	// 初始化数据库
	db := utils.IgnoreErr(gorm.Open(sqlite.Open(filepath.Join(global.DataDir, "koe.db")),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)}))
	// 迁移 schema
	db.AutoMigrate(&domain.WorkDomain{})
	db.AutoMigrate(&domain.Settings{})
	return db
}

func loadSettingsFromDB() *domain.Settings {
	settings := &domain.Settings{ID: 1, ScanDirs: []byte("[]")}
	global.GetServiceContext().DB.FirstOrCreate(settings)
	return settings
}

func web() {
	serve := router.GetGinServe()
	serve.NoRoute(func(context *gin.Context) {
		context.FileFromFS("/"+context.Request.RequestURI, http.FS(utils.IgnoreErr(fs.Sub(dist, "dist"))))
	})
	serve.Run(":" + global.GetServiceContext().Config.FlagConfig.Port)
}
