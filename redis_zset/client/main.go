package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const zsetKey = "zset:order.delayed"

func main() {
	ctx := context.Background()

	// 1. 初始化標準 Redis 連線
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// 2. 排入 3 筆「延遲任務」：score 用未來要執行的時間戳，member 是任務內容。
	// 💡 這是 ZSet 最經典的延遲佇列用法——List／Stream 都是「先進先出」，
	// 沒辦法讓某筆任務插隊排到 5 秒後才處理；ZSet 用 score 排序，
	// 天生就適合「按時間／優先度排隊」這種場景。
	fmt.Println("🎬 開始排入延遲任務...")

	aDelays := []time.Duration{6 * time.Second, 2 * time.Second, 4 * time.Second}

	for i, oDelay := range aDelays {
		sMember := fmt.Sprintf(`{"order_id": %d}`, i+1)
		oReadyTime := time.Now().Add(oDelay)
		iReadyAt := oReadyTime.Unix()
		fReadyAt := float64(iReadyAt)

		oResult := rdb.ZAdd(ctx, zsetKey, redis.Z{Score: fReadyAt, Member: sMember})
		err := oResult.Err()
		if err != nil {
			fmt.Printf("⚠️ 排入失敗: %v\n", err)
			continue
		}

		fmt.Printf("📤 [client] 已排入第 %d 筆訂單，%v 後到期\n", i+1, oDelay)
	}

	fmt.Println("✅ 範例結束（注意：排入順序是 6s→2s→4s，但到期順序會是 2s→4s→6s）")
}
