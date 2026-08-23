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

const listKey = "list:order.created"

func main() {
	// 收到中斷/終止訊號時 ctx 會被取消，下面的消費迴圈才會跟著結束，
	// 不是靠 process 被系統強制殺掉才停止。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. 初始化標準 Redis 連線
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	fmt.Printf("🚪 [server] 開始消費清單 %s，等待任務中（Ctrl+C 結束）...\n", listKey)

	// 2. 用 BRPOP 阻塞式地從清單尾端拿任務，跟 pub/sub 的 Subscribe 不同：
	// 💡 List 是「排隊等消費者上線再拿」的語意，任務會先落地存在 Redis 裡，
	// 就算這支 server 現在沒開著，client 照樣能塞任務進去，之後 server 一啟動
	// 就能立刻拿到之前堆積的任務，不會像 pub/sub 那樣直接遺失。
	for {
		aResult, err := rdb.BRPop(ctx, 5*time.Second, listKey).Result()
		if err != nil {
			if err == redis.Nil {
				// 5 秒內沒有新任務，繼續下一輪等待
				continue
			}
			if ctx.Err() != nil {
				fmt.Println("👋 [server] 收到中斷訊號，結束消費")
				return
			}
			fmt.Printf("⚠️ 消費失敗: %v\n", err)
			continue
		}

		// aResult[0] 是清單名稱，aResult[1] 是實際內容
		fmt.Printf("📩 [server] 從 %s 拿到任務：%s\n", aResult[0], aResult[1])
	}
}
