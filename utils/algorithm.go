package utils

func Max(more ...int) int {
	maxNum := more[0]
	for _, elem := range more {
		if maxNum < elem {
			maxNum = elem
		}
	}
	return maxNum
}

// Longest 最长公共子序列 https://stackoverflow.com/questions/20145795/go-longest-common-subsequence-to-print-result-array
func Longest(str1, str2 string) int {
	len1 := len(str1)
	len2 := len(str2)

	tab := make([][]int, len1+1)
	for i := range tab {
		tab[i] = make([]int, len2+1)
	}

	i, j := 0, 0
	for i = 0; i <= len1; i++ {
		for j = 0; j <= len2; j++ {
			if i == 0 || j == 0 {
				tab[i][j] = 0
			} else if str1[i-1] == str2[j-1] {
				tab[i][j] = tab[i-1][j-1] + 1
			} else {
				tab[i][j] = Max(tab[i-1][j], tab[i][j-1])
			}
		}
	}

	return tab[len1][len2]
}
