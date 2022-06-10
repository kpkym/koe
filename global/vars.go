package global

import (
	"os"
	"path/filepath"
)

var (
	DataDir = filepath.Join(func() string {
		wd, _ := os.Getwd()
		return wd
	}(), "data")
)

func init() {
	os.MkdirAll(DataDir, os.ModePerm)
}
