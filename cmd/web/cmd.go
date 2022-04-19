package web

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kpk-koe/global"
	"kpk-koe/model/domain"
	"kpk-koe/router"
	"kpk-koe/utils"
	"net/http"
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

func init() {
	Cmd.Flags().StringVarP(&global.NasCacheFile, "nasFile", "n", utils.IgnoreErr(homedir.Expand("~/config/koe/nas_koe_data/naskoe_001.json")), "NAS缓存文件")
	Cmd.Flags().StringVarP(&global.SqliteDataFile, "db", "d", utils.IgnoreErr(homedir.Expand("~/GolandProjects/kpk-koe/data/koe.db")), "sqlite数据库文件")
	Cmd.Flags().StringVarP(&global.ScanDir, "koe", "k", utils.IgnoreErr(homedir.Expand("~/ooo/koe/scan")), "扫描文件夹")
}

func web() {
	// 初始化数据库
	db := utils.IgnoreErr(gorm.Open(sqlite.Open(global.SqliteDataFile), &gorm.Config{}))
	// 迁移 schema
	db.AutoMigrate(&domain.DlDomain{})
	db.AutoMigrate(&domain.TrackDomain{})
	db.AutoMigrate(&domain.FileSystemTreeDomain{})

	global.SetServiceContext(db)

	serve := router.GetGinServe()
	serve.StaticFS("/static", http.Dir(global.ScanDir))
	serve.Run(":8081") // 监听并在 0.0.0.0:8080 上启动服务
}
