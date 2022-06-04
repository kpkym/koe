package cache

import (
	"errors"
	"github.com/allegro/bigcache/v3"
	"github.com/golang/protobuf/proto"
	"github.com/kpkym/koe/utils"
)

var (
	cache = func() *bigcache.BigCache {
		config := bigcache.DefaultConfig(100000 * 100000 * 60)
		// 缓存不过期
		config.CleanWindow = 0
		return utils.IgnoreErr(bigcache.NewBigCache(config))
	}()
)

func Set[T proto.Message](key string, t T) {
	cache.Set(key, utils.PBMarshal(t))
}

func Get[T proto.Message](key string, resp T) error {
	entry, err := cache.Get(key)

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		return err
	}

	utils.PBUnmarshal(entry, resp)
	return nil
}
