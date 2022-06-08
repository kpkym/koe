package utils

import (
	"github.com/kpkym/koe/model/pb"
	"testing"
)

func TestPBUnmarshal(t *testing.T) {
	nodeFile := &pb.PBNode{
		Type:  "file",
		Title: "123",
	}
	nodeParent := &pb.PBNode{
		Type:     "folder",
		Title:    "abc",
		Children: []*pb.PBNode{nodeFile},
	}

	PBMarshal[*pb.PBNode](nodeParent)
}

func TestPBMarshal(t *testing.T) {
	nodeFile := &pb.PBNode{
		Type:  "file",
		Title: "123",
	}
	nodeParent := &pb.PBNode{
		Type:     "folder",
		Title:    "abc",
		Children: []*pb.PBNode{nodeFile},
	}

	resp := pb.PBNode{}
	PBUnmarshal(PBMarshal[*pb.PBNode](nodeParent), &resp)
}
