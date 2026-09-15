package main

import "fmt"

type User struct {
	Name    string
	Numbers []int
}

func main() {
	//  0  1  2  3  4  5  6  7  8  9
	aNumbers1 := []int(nil)

	aNumbers2 := []int{}

	fmt.Println("aNumbers1=", aNumbers1)
	fmt.Println("aNumbers2=", aNumbers2)
	fmt.Println("==========================")

	fmt.Println("aNumbers1==nil", aNumbers1 == nil)
	fmt.Println("aNumbers2==nil", aNumbers2 == nil)

}
