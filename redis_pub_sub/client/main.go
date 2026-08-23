package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const pubSubChannel = "channel:order.created"

func main() {
	ctx := context.Background()

	// 1. 初始化標準 Redis 連線
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// 2. 模擬陸續發生 3 筆訂單建立事件
	fmt.Println("🎬 開始發布訊息...")

	for i := 1; i <= 3; i++ {
		sPayload := fmt.Sprintf(`{"order_id": %d}`, i)

		iReceivers, err := rdb.Publish(ctx, pubSubChannel, sPayload).Result()
		if err != nil {
			fmt.Printf("⚠️ 發布失敗: %v\n", err)
			continue
		}

		// 💡 Publish 回傳值是「當下有幾個訂閱者收到」。Redis Pub/Sub 是 fire-and-forget，
		// 訊息不會落地保存，發布當下沒有訂閱者在線就直接遺失，跟 Stream／隊列那種
		// 「訊息會排隊等消費者上線再拿」的語意完全不同。先把 server/main.go 跑起來，
		// 再跑這支 client，才能看到 server 端收到訊息。
		fmt.Printf("📤 [client] 已發布第 %d 筆訂單，%d 個訂閱者收到\n", i, iReceivers)

		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("✅ 範例結束")
}
