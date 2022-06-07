package utils

import "fmt"

func Str2Any(item string) any {
	return item
}

func Any2Str(item any) string {
	return fmt.Sprint(item)
}
