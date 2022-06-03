package cache

import (
	"github.com/allegro/bigcache/v3"
	"github.com/golang/protobuf/proto"
	"github.com/kpkym/koe/utils"
	"time"
)

var (
	cache = utils.IgnoreErr(bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute)))
)

func Set[T proto.Message](key string, t T) {
	cache.Set(key, utils.IgnoreErr(proto.Marshal(t)))
}

func Get[T proto.Message](key string, resp T) {
	entry, _ := cache.Get(key)
	_ = proto.Unmarshal(entry, resp)
}
