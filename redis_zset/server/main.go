package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

const zsetKey = "zset:order.delayed"

func main() {
	// 收到中斷/終止訊號時 ctx 會被取消，下面的輪詢迴圈才會跟著結束，
	// 不是靠 process 被系統強制殺掉才停止。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. 初始化標準 Redis 連線
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	fmt.Printf("🚪 [server] 開始輪詢到期任務（Ctrl+C 結束）...\n")

	oTicker := time.NewTicker(time.Second)
	defer oTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("👋 [server] 收到中斷訊號，結束輪詢")
			return
		case <-oTicker.C:
			consumeDueTasks(ctx, rdb)
		}
	}
}

// consumeDueTasks 撈出 score（到期時間戳）小於等於現在的任務，逐筆處理並移除。
// 💡 ZRangeByScore 撈出來到真的 ZRem 掉這段中間不是原子操作，兩個 server 同時搶
// 同一筆任務理論上可能重複拿到；正式場景通常會包一段 Lua script（GET 完立刻用同一次
// 呼叫 REM）來保證原子性，這裡為了範例單純用兩步驟示範概念就好。
func consumeDueTasks(ctx context.Context, rdb *redis.Client) {
	oNow := time.Now()
	iNow := oNow.Unix()
	sNow := fmt.Sprintf("%d", iNow)

	oRangeCmd := rdb.ZRangeByScore(ctx, zsetKey, &redis.ZRangeBy{
		Min: "-inf",
		Max: sNow,
	})
	aMembers, err := oRangeCmd.Result()
	if err != nil {
		fmt.Printf("⚠️ 撈取到期任務失敗: %v\n", err)
		return
	}

	for _, sMember := range aMembers {
		// 拿到就先移除，避免下一輪輪詢又撈到同一筆
		oRemCmd := rdb.ZRem(ctx, zsetKey, sMember)
		err := oRemCmd.Err()
		if err != nil {
			fmt.Printf("⚠️ 移除任務失敗: %v\n", err)
			continue
		}

		fmt.Printf("📩 [server] 任務到期，開始處理：%s\n", sMember)
	}
}
