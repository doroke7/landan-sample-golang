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
	  3. 永遠 init 的 slice 是 cap 剛好是 len
	  4. slices.Clone

	*/
	aNumbers0 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	aNumbers1 := append([]int(nil), aNumbers0...)

	aNumbers2 := make([]int, 0, 10)

	aNumbers3 := aNumbers0[:2:2]

	aNumbers4 := slices.Clone(aNumbers0) // 淺拷貝

	fmt.Println("aNumbers1=", aNumbers1)
	fmt.Println("aNumbers2=", aNumbers2)
	fmt.Println("aNumbers3=", aNumbers3)
	fmt.Println("aNumbers4=", aNumbers4)

}
