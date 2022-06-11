package cache

import "reflect"

var (
	mapCacheMap = make(map[any]any)
)

type mapCache[K comparable, T any] struct{}

func NewMapCache[K comparable, T any]() mapCache[K, T] {
	return mapCache[K, T]{}
}

func (m mapCache[K, T]) Get(key K) (value T, ok bool) {
	if v, ok := mapCacheMap[key]; ok {
		return reflect.ValueOf(v).Interface().(T), true
	}
	// logError()
	return *new(T), false
}

func (m mapCache[K, T]) Set(key K, value T) {
	mapCacheMap[key] = value
}

func (m mapCache[K, T]) GetOrSet(key K, fn func() (value T)) (value T) {
	if v, ok := m.Get(key); ok {
		return v
	}
	v := fn()
	m.Set(key, v)
	return v
}
