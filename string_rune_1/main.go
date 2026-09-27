package main

import (
	"fmt"
)

/*
Go 的單引號是 rune 格式，而且只能放一個字元
如果你做 單引號 相加 不是字串，是 整數
*/
func main() {
	sS := 'G' + 'o' + '語' + '言'

	fmt.Printf("數值: %v\n", sS) // 印出相加後的整數總和
	fmt.Printf("型別: %T\n", sS) // 印出變數型別
	fmt.Printf("字元: %c\n", sS) // 嘗試以 Unicode 字元印出
}
