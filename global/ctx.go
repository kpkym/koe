package global

import (
	"github.com/kpkym/koe/model/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

var (
	ctx serviceContext
)

type serviceContext struct {
	Config   *Config
	DB       *gorm.DB
	Settings *domain.Settings
}

func AddConfig(c *Config) {
	ctx.Config = c
}

func AddDB(db *gorm.DB) {
	ctx.DB = db
}

func AddSettings(settings *domain.Settings) {
	ctx.Settings = settings
}

func GetServiceContext() serviceContext {
	if ctx.Config == nil {
		logrus.Error("没有初始化配置类")
		os.Exit(1)
	}
	return ctx
}
