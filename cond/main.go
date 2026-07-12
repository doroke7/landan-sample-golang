package main

import (
	"fmt"
	"sync"
	"time"
)

/**
   注意： cond.Wait()  的代碼 剛執行的時候 隱含 “鎖釋放” 的動作
                            被喚起的時候 會有一個 “搶到鎖” 的動作
*/

func main() {
	var mutex sync.Mutex
	cond := sync.NewCond(&mutex)

	// 戰術狀態：是否發動總攻擊？
	startAttack := false

	var wg sync.WaitGroup
	wg.Add(3) // 有 3 個小兵要去佔地

	// 派出 3 個小兵
	for i := 1; i <= 3; i++ {
		go func(soldierID int) {
			defer wg.Done()

			fmt.Printf("🪖 兵 %d 号：正在秘密潛入陣地...\n", soldierID)
			time.Sleep(time.Duration(soldierID*200) * time.Millisecond) // 模擬行軍時間

			// --- 準備就位，進入待命狀態 ---
			cond.L.Lock() // 鎖定狀態
			
			fmt.Printf("🪖 兵 %d 号：已成功佔領陣地！【原地待命】，等待信號彈...\n", soldierID)

			// 用 for 迴圈等待指揮官的「總攻擊」命令
			for !startAttack {
				// Wait() 會釋放鎖，讓小兵靜音休眠。
				// Wait() 會釋放鎖
				// Wait() 會釋放鎖
				// Wait() 會釋放鎖
				// Wait() 會釋放鎖

				cond.Wait() 

				// 收到廣播醒來時，會自動“重新加鎖”，確保安全。
				// 收到廣播醒來時，會自動“重新加鎖”，確保安全。
				// 收到廣播醒來時，會自動“重新加鎖”，確保安全。
				// 收到廣播醒來時，會自動“重新加鎖”，確保安全。
				// 收到廣播醒來時，會自動“重新加鎖”，確保安全。
				
			}

			// 走到這裡，代表收到命令了，且持有鎖
			fmt.Printf("🪖 兵 %d 号：收到信號彈！🔥🔥🔥 發動總攻擊！！！\n", soldierID)
			
			cond.L.Unlock() // 解鎖，讓下一個醒來的小兵也能順利處理

			// 模擬戰鬥一小段時間
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("✅ 兵 %d 号：報告指揮官，任務完成，敵方陣地已拿下！\n", soldierID)
		}(i)
	}

	// 指揮官在指揮所觀察戰局
	time.Sleep(1 * time.Second)
	fmt.Println("\n🎖️ 指揮官：看來大家都已經佔領陣地，進入待命位置了。")
	fmt.Println("🎖️ 指揮官：準備發射信號彈！")

	// 指揮官下達變更狀態（加鎖保護）
	cond.L.Lock()
	startAttack = true
	cond.L.Unlock()

	// 發射信號彈，全軍出擊！
	fmt.Println("🚀 指揮官：【發射紅色信號彈】—— 全軍總攻擊！！！\n")
	cond.Broadcast()

	// 等待所有小兵打完仗並回報
	wg.Wait()
	fmt.Println("\n🎖️ 指揮官：很好！戰役完美勝利，全體收兵！")
}