package main

import (
	"fmt"
	"slices"
)

func main() {
	/*
	  如何保證 slice 都是新的？
	  1. append 一個 zero slice
	  2. 每次都make
	  3. slices.Clone

	*/
	aNumbers1 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	aNumbers11 := append([]int(nil), aNumbers1...)

	fmt.Println("aNumbers11=", aNumbers11)

	aNumbers21 := make([]int, 0, 10)

	fmt.Println("aNumbers21=", aNumbers21)

	aNumbers31 := slices.Clone(aNumbers1) // 淺拷貝

	fmt.Println("aNumbers31=", aNumbers31)

}
