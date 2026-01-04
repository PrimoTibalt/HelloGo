package main

import (
	"unicode/utf8"
)

func Ld(actual, expected string) (distance int) {
	m, n := utf8.RuneCountInString(expected)+1, utf8.RuneCountInString(actual)+1

	arr := make([][]int, m)
	for i := range arr {
		arr[i] = make([]int, n)
	}

	for i := range arr {
		for j := range arr[i] {
			distance := calcNextCell(i, j, arr, actual, expected)
			arr[i][j] = distance
		}
	}

	return arr[m-1][n-1]
}

func calcNextCell(i, j int, arr [][]int, actual, expected string) int {
	if i == 0 && j == 0 {
		return 0
	}

	if i == 0 {
		return arr[i][j-1] + 1
	}

	if j == 0 {
		return arr[i-1][j] + 1
	}

	eR := []rune(expected)[i-1]
	aR := []rune(actual)[j-1]
	if eR == aR {
		return arr[i-1][j-1]
	}

	return min(arr[i][j-1], arr[i-1][j-1], arr[i-1][j]) + 1
}
