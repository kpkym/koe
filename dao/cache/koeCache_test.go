package cache

import (
	"fmt"
	"github.com/kpkym/koe/model/pb"
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

	Set("key", nodeParent)

	resp := pb.PBNode{}
	Get("key", &resp)

	fmt.Println(resp)
}
