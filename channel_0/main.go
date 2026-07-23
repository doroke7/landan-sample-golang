package main

import (
	"fmt"
	"time"
)

func main() {

	oChannel := make(chan int)

	go func() {
		iI := 1

		for {
			time.Sleep(10 * time.Second)

			oChannel <- iI

			iI++
		}
	}()

	for {
		fmt.Println("準備 進入 loop")

		fmt.Println("讀取 channel, 如果沒有資料就會 阻塞(wait) 在這一行...")

		select {

		case iV := <-oChannel:
			fmt.Println("讀取到數字:", iV, ".")

		}
		fmt.Println("準備 離開 loop")
		fmt.Println("")

	}

	/*
		    1. 由於輸入 channel 的資料不是連續的， 大概 1秒一個
			2. 取出 channel 的地方，

	*/
	// time.Sleep(time.Second * 20)
}
