package cache

import "reflect"

var (
	mapCacheMap = make(map[string]any)
)

type mapCache[T any] struct{}

func NewMapCache[T any]() mapCache[T] {
	return mapCache[T]{}
}

func (m mapCache[T]) Get(key string) (value T, ok bool) {
	if v, ok := mapCacheMap[key]; ok {
		return reflect.ValueOf(v).Interface().(T), true
	}
	logError()
	return *new(T), false
}

func (m mapCache[T]) Set(key string, value T) {
	mapCacheMap[key] = value
}

func (m mapCache[T]) GetOrSet(key string, fn func() (value T)) (value T) {
	if v, ok := m.Get(key); ok {
		return v
	}
	v := fn()
	m.Set(key, v)
	return v
}
