package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const listKey = "list:order.created"

func main() {
	ctx := context.Background()

	// 1. 初始化標準 Redis 連線
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// 2. 模擬陸續發生 3 筆訂單建立事件
	fmt.Println("🎬 開始推送任務...")

	for i := 1; i <= 3; i++ {
		sPayload := fmt.Sprintf(`{"order_id": %d}`, i)

		// 💡 LPush 把任務推進清單的頭端，server 用 BRPop 從尾端拿，
		// 兩邊搭配起來就是先進先出（FIFO）的任務佇列。
		if err := rdb.LPush(ctx, listKey, sPayload).Err(); err != nil {
			fmt.Printf("⚠️ 推送失敗: %v\n", err)
			continue
		}

		fmt.Printf("📤 [client] 已推送第 %d 筆訂單\n", i)

		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("✅ 範例結束（就算 server 這時候才啟動，也還是拿得到這些任務）")
}
