package main

import "fmt"

type User struct {
	Id int
}

func main() {
	var oUser *User

	var oMap map[string]int = nil

	var iChNumber chan int

	var fn func()

	var x interface{}

	var aNumbers []int = nil // nil slice 會隱含 空array
	// array 是 不可nil 的 。譬如 var a [3]int

	var a [3]int // 是 0,0,0

	_ = a

	fmt.Println("oUser=", oUser)
	fmt.Println("oMap=", oMap)
	fmt.Println("aNumbers=", aNumbers)
	fmt.Println("iChNumber=", iChNumber)
	fmt.Println("fn=", fn)
	fmt.Println("x=", x)

}
