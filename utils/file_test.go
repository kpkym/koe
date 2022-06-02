package utils

import (
	"fmt"
	"testing"
)

type NameFile struct {
	A string `json:"a"`
}

func TestBuildTree(t *testing.T) {
	tree := BuildTree(false)

	fmt.Println(Marshal(tree))
}

func TestStruct(t *testing.T) {
	a := []*NameFile{}

	a = append(a, &NameFile{A: "123"})

	for _, e := range a {
		e.A = "432"
	}

	fmt.Println(*a[0])
}
