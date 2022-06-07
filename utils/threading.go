package utils

import (
	"github.com/sirupsen/logrus"
)

func GoSafe(fn func(), errMsg ...string) {
	defer func() {
		if p := recover(); p != nil {
			logrus.Error(errMsg)
		}
	}()
	go fn()
}
