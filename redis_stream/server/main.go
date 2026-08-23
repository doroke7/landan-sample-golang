package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	streamKey    = "stream:order.created"
	groupName    = "order-service"
	consumerName = "worker-1"
)

func main() {
	// 收到中斷/終止訊號時 ctx 會被取消，下面的消費迴圈才會跟著結束，
	// 不是靠 process 被系統強制殺掉才停止。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. 初始化標準 Redis 連線
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	// 2. 建立 consumer group，"$" 代表只消費「這個 group 建立之後」新進來的訊息。
	// 💡 跟 List 不同，Stream 上可以同時掛好幾個 group，每個 group 各自獨立追蹤自己
	// 消費到哪，同一則訊息可以被多個 group 各自處理一次——這是 List／pub-sub 都做不到的。
	err := rdb.XGroupCreateMkStream(ctx, streamKey, groupName, "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		fmt.Printf("⚠️ 建立 consumer group 失敗: %v\n", err)
		return
	}

	fmt.Printf("🚪 [server] 已加入 group %s，開始消費 %s（Ctrl+C 結束）...\n", groupName, streamKey)

	for {
		if ctx.Err() != nil {
			fmt.Println("👋 [server] 收到中斷訊號，結束消費")
			return
		}

		// 3. ">" 代表只要「還沒分派給任何 consumer」的新訊息
		aStreams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamKey, ">"},
			Count:    10,
			Block:    5 * time.Second,
		}).Result()
		if err != nil {
			if err == redis.Nil {
				// 5 秒內沒有新訊息，繼續下一輪等待
				continue
			}
			if ctx.Err() != nil {
				continue
			}
			fmt.Printf("⚠️ 消費失敗: %v\n", err)
			continue
		}

		for _, oStream := range aStreams {
			for _, oMsg := range oStream.Messages {
				fmt.Printf("📩 [server] 收到訊息 %s：%v\n", oMsg.ID, oMsg.Values)

				// 4. 處理完要自己 XAck，group 才知道這則訊息不用再重派給別的 consumer；
				// 沒 ack 的訊息會留在 PEL（pending entries list），可以之後用
				// XPending／XClaim 找回來重新處理，這是 Stream 特有的「至少一次」保障。
				if err := rdb.XAck(ctx, streamKey, groupName, oMsg.ID).Err(); err != nil {
					fmt.Printf("⚠️ ack 失敗: %v\n", err)
				}
			}
		}
	}
}
