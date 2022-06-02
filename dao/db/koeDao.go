package db

import (
	"errors"
	"github.com/kpkym/koe/global"
	"gorm.io/gorm"
)

type koeDB[T any] struct {
	db *gorm.DB
}

func NewKoeDB[T any]() koeDB[T] {
	return koeDB[T]{
		db: global.GetServiceContext().DB,
	}
}

func (koe koeDB[T]) GetData(model *T, code string, fn func() T) error {
	err := koe.db.First(model, "code = ?", code).Error
	// 找不到 有缓存方法
	if errors.Is(err, gorm.ErrRecordNotFound) && fn != nil {
		insertModel := fn()
		koe.db.Create(insertModel)
		*model = insertModel
		err = nil
	}
	return err
}
