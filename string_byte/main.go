package main

import "fmt"

/*
重點 ：

	go 預設是用 byte 存字串，但是是變動長度。譬如 a b c d 是一個 byte 但是 中文就是 3個 byte
*/
func main() {
	fmt.Println("")
	fmt.Println("A========================")

	sString := "Go語言"
	for index, r := range sString {
		fmt.Printf("%d: %c (Unicode: %U)\n", index, r, r)
	}
	/* 輸出:
	   0: G (Unicode: U+0047)
	   1: o (Unicode: U+006F)
	   2: 語 (Unicode: U+8A9E)  <-- 注意 index 跳到了 5，因為「語」佔了 3 bytes
	   5: 言 (Unicode: U+8A00)
	*/
	fmt.Println("")

	fmt.Println("B========================")
	for i := 0; i < len(sString); i++ {
		fmt.Print(sString[i], " ")
	}
	fmt.Println("")

	fmt.Println("C========================")

	for i := 0; i < len(sString); i++ {
		fmt.Print(string(sString[i]), " ")
	}
	fmt.Println("")

	fmt.Println("D========================")

	sRune := []rune("Go語言")

	for iIndex := 0; iIndex < len(sRune); iIndex++ {
		fmt.Print(sRune[iIndex], " ")
	}
	fmt.Println("")

	fmt.Println("E========================")

	for iIndex := 0; iIndex < len(sRune); iIndex++ {
		fmt.Print(string(sRune[iIndex]), " ")
	}
	fmt.Println("")

	fmt.Println("F========================")

}
