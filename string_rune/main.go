package main

import (
	"fmt"
)

func main() {
	// sBytes 是一個原始的 []byte 切片
	sString := "Go語言"
	sBytes := []byte("Go語言")
	sRunes := []rune("Go語言")
	/*
	 重點： go 使用 range 去操作 string , 他會隱性地 轉成 []rune 再操作
	 換言之 	sRunes := []rune("Go語言") 也是一樣的

	*/
	fmt.Println("")
	fmt.Println("A========================")
	for index, r := range sString {
		// %T 會印出變數的型別 (Type)
		fmt.Printf("Index: %d, Value: %c, Type of r: %T\n", index, r, r)
	}
	fmt.Println("")
	fmt.Println("B========================")
	for index, r := range sBytes {
		// %T 會印出變數的型別 (Type)
		fmt.Printf("Index: %d, Value: %c, Type of r: %T\n", index, r, r)
	}

	fmt.Println("")
	fmt.Println("C========================")

	for index, r := range sRunes {
		// %T 會印出變數的型別 (Type)
		fmt.Printf("Index: %d, Value: %c, Type of r: %T\n", index, r, r)
	}
}
