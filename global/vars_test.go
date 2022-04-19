package global

import (
	"fmt"
	"testing"
)

func TestGetNextNasCacheFile(t *testing.T) {
	fmt.Println(GetNextNasCacheFile())
}

func TestName(t *testing.T) {
	a := []int{1, 2, 3}
	f := true

	for i := range a {
		fmt.Println(i)

		a = append(a, 123)

		if !f {
			break
		}
	}
}
