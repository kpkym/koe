package cache

import (
	"github.com/allegro/bigcache/v3"
	"github.com/golang/protobuf/proto"
	"github.com/kpkym/koe/utils"
	"reflect"
)

var (
	cache = func() *bigcache.BigCache {
		config := bigcache.DefaultConfig(100000 * 100000 * 60)
		// 缓存不过期
		config.CleanWindow = 0
		return utils.IgnoreErr(bigcache.NewBigCache(config))
	}()
)

type bigCache[T proto.Message] struct{}

// NewBigCache 不使用
func NewBigCache[T proto.Message]() bigCache[T] { return bigCache[T]{} }

func (b bigCache[T]) Get(key string) (value T, ok bool) {
	if b, err := cache.Get(key); err == nil {
		v := new(T)
		utils.PBUnmarshal(b, *v)
		return reflect.ValueOf(v).Interface().(T), true
	}
	logError()
	return *new(T), false
}

func (b bigCache[T]) Set(key string, value T) {
	cache.Set(key, utils.PBMarshal(value))
}

func (b bigCache[T]) GetOrSet(key string, fn func() (value T)) (value T) {
	if v, ok := b.Get(key); ok {
		return v
	}
	v := fn()
	b.Set(key, v)
	return v
}
