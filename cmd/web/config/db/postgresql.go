//go:build postgresql

package db

import (
	"fmt"
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/domain"
	"github.com/kpkym/koe/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() *gorm.DB {
	config := global.GetServiceContext().Config.PostgresqlConfig

	// 初始化数据库
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Host, config.User, config.Password, config.Dbname, config.Port, config.Sslmode, config.Timezone)

	db := utils.IgnoreErr(gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}))
	// 迁移 schema
	db.AutoMigrate(&domain.WorkDomain{})
	db.AutoMigrate(&domain.Settings{})
	return db
}
