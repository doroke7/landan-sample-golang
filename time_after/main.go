package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		// 如果有兩個 chan 就用 select
		<-time.After(5 * time.Second) // 等 5 秒
		fmt.Println("hi")
	}
}
