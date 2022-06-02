package global

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"kpk-koe/cmd/web/config"
	"os"
)

var (
	ctx serviceContext
)

type serviceContext struct {
	Config *config.Config
	DB     *gorm.DB
}

func SetServiceContext(c *config.Config) {
	ctx = serviceContext{Config: c}
}

func AddDB(db *gorm.DB) {
	ctx.DB = db
}

func GetServiceContext() serviceContext {
	if ctx.Config == nil {
		logrus.Error("没有初始化配置类")
		os.Exit(1)
	}
	return ctx
}
