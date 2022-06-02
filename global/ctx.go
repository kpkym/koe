package global

import (
	"gorm.io/gorm"
	"kpk-koe/cmd/web/config"
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
	return ctx
}
