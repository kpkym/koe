package db

import (
	"fmt"
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

func (koe koeDB[T]) Page(model *[]T, pageRequest dto.PageRequest, fn func(*gorm.DB)) int64 {
	var count int64

	fn(koe.db)
	koe.db.Count(&count)
	koe.db.Offset((pageRequest.Page - 1) * pageRequest.Size).
		Limit(pageRequest.Size).
		Order(fmt.Sprintf("%s %s", pageRequest.Order, pageRequest.Sort)).
		Find(model)
	return count
}
