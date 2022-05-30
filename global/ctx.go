package global

import (
	"gorm.io/gorm"
	"kpk-koe/cmd/web/config"
)

var (
	ctx serviceContext
)

type serviceContext struct {
	DB     *gorm.DB
	Config *config.Config
}

func SetServiceContext(db *gorm.DB, c *config.Config) {
	ctx = serviceContext{DB: db, Config: c}
}

func GetServiceContext() serviceContext {
	return ctx
}
