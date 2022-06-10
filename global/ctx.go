package global

import (
	"github.com/kpkym/koe/cmd/web/config"
	"github.com/kpkym/koe/model/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

var (
	ctx serviceContext
)

type serviceContext struct {
	Config   *config.Config
	DB       *gorm.DB
	Settings *domain.Settings
}

func SetServiceContext(c *config.Config, db *gorm.DB) {
	ctx = serviceContext{Config: c, DB: db}
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
