package utils

import "sync/atomic"

var (
	uuid uint32 = 0
)

func NextUUID() uint32 {
	return atomic.AddUint32(&uuid, 1)
}
