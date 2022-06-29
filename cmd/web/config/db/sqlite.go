//go:build !postgresql

package db

import (
	"fmt"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"path/filepath"
)

func Init() *gorm.DB {
	// 初始化数据库
	db := utils.IgnoreErr(gorm.Open(sqlite.Open(filepath.Join(global.DataDir, global.GetServiceContext().Config.SqliteConfig.Filename)),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)}))
	// 迁移 schema
	db.AutoMigrate(&domain.WorkDomain{})
	db.AutoMigrate(&domain.Settings{})
	return db
}

func SelectGroupBy(g *gorm.DB, category string) {
	g.Select("j.value name, COUNT(1) count").
		Joins(fmt.Sprintf("join json_each(%s) j", category)).
		Group("j.value").Order("count DESC")
}
