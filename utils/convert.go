package utils

import (
	"fmt"
	"strconv"
)
import "golang.org/x/exp/constraints"

func Str2Any(item string) any {
	return item
}

func Any2Str(item any) string {
	return fmt.Sprint(item)
}

func Str2Num[T constraints.Integer](item string) T {
	return T(IgnoreErr(strconv.ParseInt(item, 10, 64)))
}
