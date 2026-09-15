package main

import "fmt"

func main() {
	aNumbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("len=%d cap=%d %v\n", len(aNumbers), cap(aNumbers), aNumbers)

	aNumbers = append(aNumbers, 10)
	fmt.Printf("len=%d cap=%d %v\n", len(aNumbers), cap(aNumbers), aNumbers)

	aNumbers = append(aNumbers, 11)
	fmt.Printf("len=%d cap=%d %v\n", len(aNumbers), cap(aNumbers), aNumbers)

	aNumbers = append(aNumbers, 12)
	fmt.Printf("len=%d cap=%d %v\n", len(aNumbers), cap(aNumbers), aNumbers)

	aNumbers = append(aNumbers, 13)
	fmt.Printf("len=%d cap=%d %v\n", len(aNumbers), cap(aNumbers), aNumbers)

	aNumbers = append(aNumbers, 14)
	fmt.Printf("len=%d cap=%d %v\n", len(aNumbers), cap(aNumbers), aNumbers)

	aNumbers = append(aNumbers, 15)
	fmt.Printf("len=%d cap=%d %v\n", len(aNumbers), cap(aNumbers), aNumbers)
}
