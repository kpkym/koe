package db

import (
	"github.com/kpkym/koe/global"
	"github.com/kpkym/koe/model/dto"
	"gorm.io/gorm"
)

type koeDB[T any] struct {
	db *gorm.DB
}

func NewKoeDB[T any](t any) koeDB[T] {
	return koeDB[T]{
		db: global.GetServiceContext().DB.Model(t),
	}
}

func (koe koeDB[T]) ListByCode(model *T, codes []string) error {
	return koe.db.Where("code IN ?", codes).Find(model).Error
}

func (koe koeDB[T]) Page(model *[]T, pageRequest dto.PageRequest) int64 {
	var count int64
	koe.db.Count(&count)                                                                         // 总行数
	koe.db.Offset((pageRequest.Page - 1) * pageRequest.Size).Limit(pageRequest.Size).Find(model) // 查询pageindex页的数据
	return count
}
