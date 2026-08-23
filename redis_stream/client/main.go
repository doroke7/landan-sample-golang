package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const streamKey = "stream:order.created"

func main() {
	ctx := context.Background()

	// 1. 初始化標準 Redis 連線
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// 2. 模擬陸續發生 3 筆訂單建立事件
	fmt.Println("🎬 開始寫入訊息...")

	for i := 1; i <= 3; i++ {
		// 💡 XAdd 的 ID 給 "*" 讓 Redis 自動產生（時間戳-序號），每則訊息在
		// stream 上都有唯一、遞增的 ID，天生就有序，這點跟 List／pub-sub 又不一樣。
		sId, err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: streamKey,
			ID:     "*",
			Values: map[string]any{
				"order_id": i,
			},
		}).Result()
		if err != nil {
			fmt.Printf("⚠️ 寫入失敗: %v\n", err)
			continue
		}

		fmt.Printf("📤 [client] 已寫入第 %d 筆訂單，訊息 ID：%s\n", i, sId)

		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("✅ 範例結束（訊息會留在 stream 上，server 晚點才啟動也還是拿得到）")
}
