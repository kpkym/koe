package cache

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

type Cache[K comparable, T any] interface {
	Get(key K) (value T, ok bool)
	Set(key K, value T)
	GetOrSet(key K, fn func() (value T)) (value T)
}

func logError() {
	logrus.Error("缓存不存在: ", string(debug.Stack()))
}
