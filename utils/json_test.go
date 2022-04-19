package utils

import (
	"fmt"
	"testing"
)

type Name struct {
	A string `json:"a"`
}

func TestUnmarshal(t *testing.T) {
	s := `{"a":"abc"}`

	var n Name
	Unmarshal(s, &n)

	fmt.Println(Marshal(n))
}

func TestMarshal(t *testing.T) {
	name := Name{A: "abc"}

	marshal := Marshal(name)

	fmt.Println(marshal)
}
