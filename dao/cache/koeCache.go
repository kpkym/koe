package cache

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

type Cache[T any] interface {
	Get(key string) (value T, ok bool)
	Set(key string, value T)
	GetOrSet(key string, fn func() (value T)) (value T)
}

func logError() {
	logrus.Error("缓存不存在: ", string(debug.Stack()))
}
