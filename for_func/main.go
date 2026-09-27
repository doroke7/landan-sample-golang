package main

import (
	"fmt"
	"time"
)

// hasWork 寫在 for 上面,被下面的 for 拿來當條件用(就像 runningLauncher(pid, exe))。
func hasWork(n int) bool {
	return n <= 1000
}

// nextN 也寫在 for 上面,把遞增這件事包成 func。
func nextN(n int) int {
	return n + 1
}

func main() {
	// 變種 1:條件式的 for,只有 hasWork(n),index 的遞增(n++)寫在迴圈本體裡。
	n := 1
	for hasWork(n) {
		select {
		case <-time.After(time.Second):
			fmt.Println(n)
			n++
		}
	}

	// 變種 2:三段式 for,index(n)和遞增(n++)都寫進 for 的括號裡,條件一樣用 hasWork。
	for n := 1; hasWork(n); n++ {
		select {
		case <-time.After(time.Second):
			fmt.Println(n)
		}
	}

	// 變種 3:三段式 for,遞增改用 func(nextN),不寫 n++。
	for n := 1; hasWork(n); n = nextN(n) {
		select {
		case <-time.After(time.Second):
			fmt.Println(n)
		}
	}

	// 變種 4:無窮 for{},一樣呼叫 hasWork 判斷要不要 break,效果跟前面完全一樣。
	n = 1
	for {
		if !hasWork(n) {
			break
		}
		select {
		case <-time.After(time.Second):
			fmt.Println(n)
			n++
		}
	}

	/*
	  for 提供幾種語法
	  1. for 當作 while condition 用，有進入條件
	  2. for (i=0, func(), i++)
	  4. for 當做 while true 用，無進入條件 但是會 break 出去
	*/
}
