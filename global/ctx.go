package global

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	ctx serviceContext
)

type serviceContext struct {
	DB *gorm.DB
	V  *viper.Viper
}

func SetServiceContext(db *gorm.DB, v *viper.Viper) {
	ctx = serviceContext{DB: db, V: v}
}

func GetServiceContext() serviceContext {
	return ctx
}
