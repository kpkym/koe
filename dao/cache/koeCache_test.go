package cache

import (
	"fmt"
	"github.com/kpkym/koe/model/dto"
	"github.com/kpkym/koe/model/pb"
	"github.com/kpkym/koe/utils"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {

	nodeFile := &pb.PBNode{
		Type:  "file",
		Title: "123",
	}

	nodeParent := &pb.PBNode{
		Type:     "folder",
		Title:    "abc",
		Children: []*pb.PBNode{nodeFile},
	}

	NewBigCache[*pb.PBNode]().Set("key", nodeParent)

	resp := pb.PBNode{}
	NewBigCache[*pb.PBNode]().Get("key")

	fmt.Println(resp)
}

type Generic[T any] struct{}

func (receiver Generic[T]) fullName(t T) {

	of := utils.GetEle(reflect.ValueOf(t)).Type()
	// of := reflect.TypeOf(elem)
	logrus.Info("log [", of.Name())
	logrus.Info("log [", of.PkgPath())
	logrus.Info("log [", filepath.Join(of.PkgPath(), of.Name()))
}

func TestPrintFullName(t *testing.T) {
	Generic[*dto.PageRequest]{}.fullName(&dto.PageRequest{})
}

func TestMapCache(t *testing.T) {
	var c Cache[string]
	c = NewMapCache[string]()

	c.Set("abc", "123")
	if get, ok := c.Get("abc"); ok {
		fmt.Println(get)
	}

}
