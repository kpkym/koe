package web

import (
	"bytes"
	_ "embed"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

//go:embed config/config.toml
var configString string

func init() {
	Cmd.Flags().StringVarP(&global.NasCacheFile, "nasFile", "n", utils.IgnoreErr(homedir.Expand("~/config/koe/nas_koe_data/naskoe_001.json")), "NAS缓存文件")
	Cmd.Flags().StringVarP(&global.SqliteDataFile, "db", "d", utils.IgnoreErr(homedir.Expand("~/GolandProjects/kpk-koe/data/koe.db")), "sqlite数据库文件")
	Cmd.Flags().StringVarP(&global.ScanDir, "koe", "k", utils.IgnoreErr(homedir.Expand("~/ooo/koe/scan")), "扫描文件夹")
}

func web() {
	global.SetServiceContext(initDB(), initConfig())

	serve := router.GetGinServe()
	serve.StaticFS("/static", http.Dir(global.ScanDir))
	serve.Run(":8081")
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

func initConfig() *viper.Viper {
	buffer := new(bytes.Buffer)
	buffer.WriteString(configString)

	v := viper.New()
	v.SetConfigType("toml")
	v.ReadConfig(buffer)
	return v
}
