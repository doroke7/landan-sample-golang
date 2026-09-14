package main

import "fmt"

type User struct {
	Id int
}

/*
為什麼 slice map chan 需要用 make ，用簡體中文回答













因为 slice、map、chan 都属于比较特殊的类型：变量本身和它们真正运行所需的底层数据结构不是一回事。

*/

func main() {
	var oUser *User

	var fn func()

	var x interface{}

	var oMap map[string]int = nil
	var iChNumber chan int
	var aNumbers []int = nil // nil slice 會隱含 空array
	// array 是 不可nil 的 。譬如 var a [3]int

	var aIntegers [3]int // 是 0,0,0

	_ = aIntegers

	fmt.Println("oUser=", oUser)
	fmt.Println("oMap=", oMap)
	fmt.Println("aNumbers=", aNumbers)
	fmt.Println("iChNumber=", iChNumber)
	fmt.Println("fn=", fn)
	fmt.Println("x=", x)

}
