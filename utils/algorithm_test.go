package utils

import (
	"fmt"
	"testing"
)

func TestLongest(t *testing.T) {
	str1 := "AGGTABTABTABTAB"
	str2 := "GXTXAYBTABTABTAB"
	fmt.Println(Longest(str1, str2))
	// Actual Longest Common Subsequence: GTABTABTABTAB
	// GGGGGTAAAABBBBTTTTAAAABBBBTTTTAAAABBBBTTTTAAAABBBB
	// 13

	str3 := "AGGTABGHSRCBYJSVDWFVDVSBCBVDWFDWVV"
	str4 := "GXTXAYBRGDVCBDVCCXVXCWQRVCBDJXCVQSQQ"
	fmt.Println(Longest(str3, str4))
	// Actual Longest Common Subsequence: ?
	// GGGTTABGGGHHRCCBBBBBBYYYJSVDDDDDWWWFDDDDDVVVSSSSSBCCCBBBBBBVVVDDDDDWWWFWWWVVVVVV
	// 14
}
