package global

import "gorm.io/gorm"

var (
	ctx serviceContext
)

type serviceContext struct {
	DB *gorm.DB
}

func SetServiceContext(db *gorm.DB) {
	ctx = serviceContext{DB: db}
}

func GetServiceContext() serviceContext {
	return ctx
}
