package cache

import (
	"encoding/json"
	"errors"
	"github.com/allegro/bigcache/v3"
	"github.com/golang/protobuf/proto"
	"github.com/kpkym/koe/utils"
	"github.com/sirupsen/logrus"
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

func GetOrSet[T proto.Message](key string, resp T, fn func() T) bool {
	// 存在缓存
	if err := Get(key, resp); err == nil {
		return true
	}
	// 不存在缓存, 设置缓存
	Set[T](key, fn())
	Get(key, resp)

	return false
}

func SetJSON[T any](key string, t T) {
	cache.Set(key, utils.IgnoreErr(json.Marshal(t)))
}

func GetJSON[T any](key string) T {
	entry, err := cache.Get(key)
	resp := new(T)

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		return *resp
	}
	json.Unmarshal(entry, resp)
	return *resp
}

func GetOrSetJSON[T any](key string, fn func() T) T {
	entry, err := cache.Get(key)

	resp := new(T)

	// 存在缓存
	if err == nil {
		logrus.Info("[GetOrSetJSON]读缓存:", key)
		utils.Unmarshal(string(entry), resp)
		return *resp
	}
	result := fn()
	cache.Set(key, utils.IgnoreErr(json.Marshal(result)))
	return result
}
